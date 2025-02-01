package repo

import (
	"context"

	"go-loyalty-system/internal/entity"
	"go-loyalty-system/pkg/logging"
	"go-loyalty-system/pkg/postgres"

	"go.uber.org/zap"
)

const _defaultEntityCap = 64

type GopherMartRepo struct {
	pg     *postgres.Postgres
	Logger *logging.ZapLogger
}

func NewUserRepo(pg *postgres.Postgres, l *logging.ZapLogger) *GopherMartRepo {
	return &GopherMartRepo{
		pg:     pg,
		Logger: l,
	}
}

func (g *GopherMartRepo) GetUser(ctx context.Context) ([]entity.User, error) {
	sql, _, err := g.pg.Builder.
		Select("login, email").
		From("users").
		ToSql()
	if err != nil {
		g.Logger.ErrorCtx(ctx, "TranslationRepo - GetUser - r.Builder: %w", zap.Error(err))
		return nil, err
	}

	rows, err := g.pg.Pool.Query(ctx, sql)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "TranslationRepo - GetUser - r.Pool.Query: %w", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	entities := make([]entity.User, 0, _defaultEntityCap)

	for rows.Next() {
		e := entity.User{}

		err = rows.Scan(&e.Login, &e.Email)
		if err != nil {
			g.Logger.ErrorCtx(ctx, "TranslationRepo - GetUser - rows.Scan: %w", zap.Error(err))
			return nil, err
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (g *GopherMartRepo) RegisterUser(ctx context.Context, u entity.User) error {
	sql, args, err := g.pg.Builder.
		Insert("users").
		Columns("login, email").
		Values(u.Login, u.Email).
		ToSql()
	if err != nil {
		g.Logger.ErrorCtx(ctx, "TranslationRepo - RegisterUser - r.Builder: %w", zap.Error(err))
		return err
	}

	_, err = g.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "TranslationRepo - RegisterUser - r.Pool.Exec: %w", zap.Error(err))
		return err
	}

	return nil
}
