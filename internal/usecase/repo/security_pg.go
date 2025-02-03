package repo

import (
	"context"
	"go-loyalty-system/internal/entity"

	"go.uber.org/zap"
)

func (g *GopherMartRepo) CreateToken(ctx context.Context, t *entity.Token) error {

	sql, args, err := g.pg.Builder.
		Insert("token").
		Columns("user_id, creation_date, used_at").
		Values(t.UserID, t.UsedAt, t.UsedAt).
		ToSql()
	if err != nil {
		g.Logger.ErrorCtx(ctx, "TranslationRepo - CreateToken - r.Builder: %w", zap.Error(err))
		return err
	}

	_, err = g.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "TranslationRepo - CreateToken - r.Pool.Exec: %w", zap.Error(err))
		return err
	}

	return nil
}
