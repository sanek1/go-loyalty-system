package repo

import (
	"context"
)

func (g *GopherMartRepo) SaveAccrual(ctx context.Context, orderNumber string, status string, accrual float32) error {
	// get order data by order number
	// order,err:=g.GetOrderByNumber(ctx, orderNumber)
	// if err != nil {
	// 	g.Logger.ErrorCtx(ctx, "SaveAccrual: %w", zap.Error(err))
	// 	return fmt.Errorf("SaveAccrual: %w", err)
	// }

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

	//time.Sleep(20 * time.Second)

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
