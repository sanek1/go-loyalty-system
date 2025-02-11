package repo

import (
	"context"
	"go-loyalty-system/internal/entity"
	"time"
)

func (g *GopherMartRepo) SetOrders(ctx context.Context, userID uint, o entity.Order) error {
	user, err := g.GetUserByID(ctx, entity.User{ID: userID})
	if err != nil {
		return g.logAndReturnError(ctx, "SetOrders - r.GetUser", err)
	}

	sql, args, err := g.pg.Builder.
		Insert("orders").
		Columns("user_id", "status_id, creation_date", "number").
		Values(user.ID, 1, time.Now(), o.Number).
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
