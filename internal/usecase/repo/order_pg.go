package repo

import (
	"context"
	"errors"
	"fmt"
	"go-loyalty-system/internal/entity"
	"time"

	"github.com/jackc/pgx"
)

func (g *GopherMartRepo) SetOrders(ctx context.Context, userID uint, o entity.Order) error {
	if err := g.validateOrder(o); err != nil {
        return fmt.Errorf("order validation failed: %w", err)
    }


	// user, err := g.GetUserByID(ctx, entity.User{ID: userID})
	// if err != nil {
	// 	return g.logAndReturnError(ctx, "SetOrders - r.GetUser", err)
	// }

	sql, args, err := g.pg.Builder.
		Insert("orders").
		Columns("user_id", "status_id, creation_date", "uploaded_at", "number").
		Values(userID, entity.OrderStatusNewID, time.Now(), time.Now(), o.Number).
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
            id,
            status_id,
            uploaded_at
        FROM orders
        WHERE user_id = $1
        ORDER BY uploaded_at DESC
    `

	rows, err := g.pg.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GetUserOrders - Query", err)
	}
	defer rows.Close()

	var orders []entity.OrderResponse
	for rows.Next() {
		var order entity.OrderResponse
		//var accrual sql.NullFloat64

		err := rows.Scan(
			&order.Number,
			&order.StatusID,
			//&accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return nil, g.logAndReturnError(ctx, "GopherMartRepo - GetUserOrders - Scan", err)
		}

		// if accrual.Valid {
		// 	order.Accrual = &accrual.Float64
		// }

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, g.logAndReturnError(ctx, "GopherMartRepo - GetUserOrders - rows.Err", err)
	}

	return orders, nil
}

func (g *GopherMartRepo) validateOrder(order entity.Order) error {
	if order.Number == "" {
		return entity.ErrInvalidOrder
	}

	exists, err := g.OrderExists(context.Background(), order.Number)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if exists {
		return entity.ErrOrderExists
	}

	return nil
}

func (r *GopherMartRepo) OrderExists(ctx context.Context, orderNumber string) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1 
            FROM orders 
            WHERE number = $1
        )
    `
	var exists bool

	err := r.pg.Pool.QueryRow(ctx, query, orderNumber).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("GopherMartRepo - OrderExists - QueryRow: %w", err)
	}

	return exists, nil
}
