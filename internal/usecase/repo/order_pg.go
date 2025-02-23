package repo

import (
	"context"
	"errors"
	"fmt"
	"go-loyalty-system/internal/entity"

	"github.com/jackc/pgx"
)

const (
	queryUserOrders = `
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
	queryCheckOrderExist = `
		SELECT user_id 
		FROM orders 
		WHERE number = $1 
		LIMIT 1
	`
	queryGetOrderByNumber = `
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
	queryOrderExists = `
	SELECT EXISTS (
	SELECT 1 
	FROM orders 
	WHERE number = $1
	)
	`
	querySetOrders = `
	INSERT INTO orders (user_id, status_id, creation_date, uploaded_at, number) 
	VALUES ($1, $2, $3, $4, $5)
	`
)

func (g *GopherMartRepo) SetOrders(ctx context.Context, userID uint, o entity.Order) error {
	_, err := g.pool.Exec(ctx, querySetOrders, userID, entity.OrderStatusNewID, o.CreatedAt, o.UploadedAt, o.Number)
	if err != nil {
		return g.logAndReturnError(ctx, "SetOrders - Exec", err)
	}
	return nil
}

func (g *GopherMartRepo) GetUserOrders(ctx context.Context, userID uint) ([]entity.OrderResponse, error) {
	rows, err := g.pg.Pool.Query(ctx, queryUserOrders, userID)
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

func (g *GopherMartRepo) CheckOrderExistence(ctx context.Context,
	orderNumber string,
	userID uint) (orderExists bool, existingUserID uint, err error) {
	err = g.pg.Pool.QueryRow(ctx, queryCheckOrderExist, orderNumber).Scan(&existingUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, 0, g.logAndReturnError(ctx, "Order does not exist", err)
		}
		return false, 0, nil
	}
	return true, existingUserID, nil
}

func (g *GopherMartRepo) GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.OrderResponse, error) {
	row := g.pg.Pool.QueryRow(ctx, queryGetOrderByNumber, orderNumber)
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
		return g.logAndReturnError(ctx, "Failed to check order", err)
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
	var exists bool
	err := g.pg.Pool.QueryRow(ctx, queryOrderExists, orderNumber).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("GopherMartRepo - OrderExists - QueryRow: %w", err)
	}
	return exists, nil
}
