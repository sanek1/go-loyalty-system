package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (g *GopherMartRepo) ExistOrderAccrual(ctx context.Context, orderNumber string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM accrual 
			WHERE order_id = (SELECT id FROM orders WHERE number = $1)
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

func (g *GopherMartRepo) SaveAccrual(ctx context.Context, orderNumber, status string, accrual float32) error {
	query := `
        INSERT INTO accrual (order_id, status_id, accrual)
        VALUES (
            (SELECT id FROM orders WHERE number = $1),
            (SELECT id FROM accrual_statuses WHERE status = $2),
            $3
           
        )
    
           
    `

	_, err := g.pg.Pool.Exec(ctx, query, orderNumber, status, accrual)
	if err != nil {
		return g.logAndReturnError(ctx, "GopherMartRepo -SaveAccrual -  SaveAccrual", err)
	}
	g.Logger.InfoCtx(ctx, "SaveAccrual", zap.String("orderNumber", orderNumber), zap.String("status", status), zap.Float32("accrual", accrual))

	query = `
        update orders
		set status_id = (select id from statuses where status = $2)
		where number = $1  
    `
	_, err = g.pg.Pool.Exec(ctx, query, orderNumber, status)
	if err != nil {
		return g.logAndReturnError(ctx, "GopherMartRepo -SaveAccrual -  update orders", err)
	}

	query = `
        INSERT INTO balance (user_id, current_balance, withdrawn)
        VALUES ((select user_id from orders where number = $1), $2, 0)
        
    `

	_, err = g.pool.Exec(ctx, query, orderNumber, accrual)
	if err != nil {
		return g.logAndReturnError(ctx, "GopherMartRepo -SaveAccrual -  update balance", err)
	}

	return nil
}

func (g *GopherMartRepo) GetUnprocessedOrders(ctx context.Context) ([]string, error) {
	query := `
        SELECT o.number
        FROM orders o
		left join accrual as a ON a.order_id = o.id
		left join accrual_statuses as s ON a.status_id = s.id
        WHERE a.id IS NULL
           OR s.status NOT IN ('PROCESSED', 'INVALID')
        ORDER BY o.uploaded_at ASC
    `

	rows, err := g.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GopherMartRepo -GetUnprocessedOrders - Query", err)
	}
	defer rows.Close()

	var orders []string
	for rows.Next() {
		var order string
		if err := rows.Scan(&order); err != nil {
			return nil, g.logAndReturnError(ctx, "GopherMartRepo -GetUnprocessedOrders -  QueryRow scan", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}
