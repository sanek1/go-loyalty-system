package repo

import (
	"context"
	"go-loyalty-system/internal/entity"
)

func (g *GopherMartRepo) CreateToken(ctx context.Context, t *entity.Token) error {
	sql, args, err := g.pg.Builder.
		Insert("token").
		Columns("id", "user_id, creation_date, used_at").
		Values(t.ID, t.UserID, t.CreationDate, t.UsedAt).
		ToSql()
	if err != nil {
		return g.logAndReturnError(ctx, "TranslationRepo - CreateToken - r.Builder", err)
	}

	_, err = g.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return g.logAndReturnError(ctx, "TranslationRepo - CreateToken - r.Pool.Exec", err)
	}
	return nil
}
