package repo

import (
	"context"
	"errors"
	"fmt"
	"go-loyalty-system/internal/entity"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

func (g *GopherMartRepo) SetOrders(ctx context.Context, userID uint, o entity.Order) error {
	sql, args, err := g.pg.Builder.
		Insert("orders").
		Columns("user_id", "status_id, creation_date", "uploaded_at", "number").
		Values(userID, entity.OrderStatusNewID, o.CreatedAt, o.UploadedAt, o.Number).
		ToSql()

	if err != nil {
		return g.logAndReturnError(ctx, "SetOrders - r.Builder", err)
	}

	_, err = g.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return g.logAndReturnError(ctx, "SetOrders - r.Pool.Exec", err)
	}

	return nil
}

func (g *GopherMartRepo) GetUserOrders(ctx context.Context, userID uint) ([]entity.OrderResponse, error) {
	query := `
        SELECT 
            CAST(o.number AS TEXT) as number,
            s.status,
            COALESCE(a.accrual, 0) as accrual,
            o.uploaded_at
        FROM orders as o
		left join statuses as s ON o.status_id = s.id 
		left join accrual as a ON o.id = a.order_id 
		WHERE o.user_id = $1
        ORDER BY uploaded_at DESC
    `

	rows, err := g.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GetUserOrders - Query", err)
	}
	defer rows.Close()

	var orders []entity.OrderResponse
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		var order entity.OrderResponse
		err := rows.Scan(
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return nil, g.logAndReturnError(ctx, "GopherMartRepo - GetUserOrders - Scan", err)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, g.logAndReturnError(ctx, "GopherMartRepo - GetUserOrders - rows.Err", err)
	}

	return orders, nil
}

func (g *GopherMartRepo) CheckOrderExistence(ctx context.Context, orderNumber string, userID uint) (bool, uint, error) {
	query := `
        SELECT user_id 
        FROM orders 
        WHERE number = $1 
        LIMIT 1
    `

	var existingUserID uint
	err := g.pg.Pool.QueryRow(ctx, query, orderNumber).Scan(&existingUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, 0, g.logAndReturnError(ctx, "Order does not exist", err)
		}
		return false, 0, nil
	}
	return true, existingUserID, nil
}

func (g *GopherMartRepo) GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.OrderResponse, error) {
	query := `
		SELECT 
			o.id,
			CAST(o.number AS TEXT) as number,
            s.status,
            COALESCE(a.accrual, 0) as accrual,
            o.uploaded_at
		FROM orders as o
		left join statuses as s ON o.status_id = s.id 
		left join accrual as a ON o.id = a.order_id 
		WHERE o.number = $1
	`

	row := g.pg.Pool.QueryRow(ctx, query, orderNumber)

	order := &entity.OrderResponse{}
	err := row.Scan(
		&order.ID,
		&order.Number,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
	)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GopherMartRepo - GetOrderByNumber - QueryRow", err)
	}

	return order, nil
}

func (g *GopherMartRepo) ValidateOrder(order entity.Order, userID uint) error {
	ctx := context.Background()
	if len(order.Number) < 5 || len(order.Number) > 20 {
		g.Logger.ErrorCtx(ctx, "order number length must be between 5 and 20")
		return entity.ErrInvalidOrder
	}
	if order.Number == "" {
		g.Logger.ErrorCtx(ctx, "empty order number")
		return entity.ErrInvalidOrder
	}
	for _, r := range order.Number {
		if r < '0' || r > '9' {
			g.Logger.ErrorCtx(ctx, "order number contains non-numeric characters")
			return entity.ErrInvalidOrder
		}
	}

	exists, existingUserID, err := g.CheckOrderExistence(ctx, order.Number, userID)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "Failed to check order: %w", zap.Error(err))
		return fmt.Errorf("failed to check order: %w", err)
	}

	if exists {
		if existingUserID == userID {
			return entity.ErrOrderExistsThisUser
		}
		return entity.ErrOrderExistsOtherUser
	}

	if !validateLuhn(order.Number) {
		return entity.ErrInvalidOrder
	}
	return nil
}

func validateLuhn(number string) bool {
	maxlen := 9
	sum := 0
	isSecond := false

	for i := len(number) - 1; i >= 0; i-- {
		d := int(number[i] - '0')

		if isSecond {
			d *= 2
			if d > maxlen {
				d -= maxlen
			}
		}

		sum += d
		isSecond = !isSecond
	}

	return sum%10 == 0
}

func (g *GopherMartRepo) OrderExists(ctx context.Context, orderNumber string) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1 
            FROM orders 
            WHERE number = $1
        )
    `
	var exists bool

	err := g.pg.Pool.QueryRow(ctx, query, orderNumber).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("GopherMartRepo - OrderExists - QueryRow: %w", err)
	}

	return exists, nil
}
