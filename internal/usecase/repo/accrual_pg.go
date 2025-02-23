package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const (
	queryExistOrderAccrual = `
        SELECT EXISTS (
            SELECT 1 
            FROM accrual a
            JOIN orders o ON o.id = a.order_id 
            WHERE o.number = $1
        )`

	queryInsertAccrual = `
        WITH order_data AS (
            SELECT id, user_id 
            FROM orders 
            WHERE number = $1
        )
        INSERT INTO accrual (order_id, status_id, accrual)
        SELECT 
            od.id,
            (SELECT id FROM accrual_statuses WHERE status = $2),
            $3
        FROM order_data od
        RETURNING order_id`

	queryUpdateOrderStatus = `
        UPDATE orders
        SET status_id = (SELECT id FROM statuses WHERE status = $2)
        WHERE number = $1`

	queryInsertBalance = `
        INSERT INTO balance (user_id, current_balance, withdrawn)
        SELECT user_id, $2, 0
        FROM orders
        WHERE number = $1`

	queryUnprocessedOrders = `
        SELECT o.number
        FROM orders o
        LEFT JOIN accrual a ON a.order_id = o.id
        LEFT JOIN accrual_statuses s ON a.status_id = s.id
        WHERE a.id IS NULL
           OR s.status NOT IN ('PROCESSED', 'INVALID')
        ORDER BY o.uploaded_at ASC`
)

// ExistOrderAccrual проверяет существование начисления для заказа
func (g *GopherMartRepo) ExistOrderAccrual(ctx context.Context, orderNumber string) (bool, error) {
	var exists bool
	err := g.pg.Pool.QueryRow(ctx, queryExistOrderAccrual, orderNumber).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check order accrual existence: %w", err)
	}
	return exists, nil
}

// SaveAccrual сохраняет информацию о начислении баллов
func (g *GopherMartRepo) SaveAccrual(ctx context.Context, orderNumber, status string, accrual float32) error {
	tx, err := g.pg.Pool.Begin(ctx)
	if err != nil {
		return g.logAndReturnError(ctx, "SaveAccrual - begin transaction", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Вставляем начисление
	var orderID int
	err = tx.QueryRow(ctx, queryInsertAccrual, orderNumber, status, accrual).Scan(&orderID)
	if err != nil {
		return g.logAndReturnError(ctx, "SaveAccrual - insert accrual", err)
	}

	// Обновляем статус заказа
	_, err = tx.Exec(ctx, queryUpdateOrderStatus, orderNumber, status)
	if err != nil {
		return g.logAndReturnError(ctx, "SaveAccrual - update order status", err)
	}

	// Обновляем баланс пользователя
	_, err = tx.Exec(ctx, queryInsertBalance, orderNumber, accrual)
	if err != nil {
		return g.logAndReturnError(ctx, "SaveAccrual - insert balance", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return g.logAndReturnError(ctx, "SaveAccrual - commit transaction", err)
	}

	g.Logger.InfoCtx(ctx, "accrual saved successfully",
		zap.String("orderNumber", orderNumber),
		zap.String("status", status),
		zap.Float32("accrual", accrual))

	return nil
}

// GetUnprocessedOrders возвращает список необработанных заказов
func (g *GopherMartRepo) GetUnprocessedOrders(ctx context.Context) ([]string, error) {
	rows, err := g.pg.Pool.Query(ctx, queryUnprocessedOrders)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GetUnprocessedOrders - query unprocessed orders", err)
	}
	defer rows.Close()

	var orders []string
	for rows.Next() {
		var order string
		if err := rows.Scan(&order); err != nil {
			return nil, g.logAndReturnError(ctx, "GetUnprocessedOrders - scan order number", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, g.logAndReturnError(ctx, "GetUnprocessedOrders - iterate orders", err)
	}

	return orders, nil
}
