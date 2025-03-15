package repo

import (
	"context"
	"fmt"

	"go-loyalty-system/pkg/logging"
	"go-loyalty-system/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const _defaultEntityCap = 64
const GopherMartRepoName = "GopherMartRepo"

type GopherMartRepo struct {
	pg     *postgres.Postgres
	Logger *logging.ZapLogger
	pool   *pgxpool.Pool
}

// func NewUserRepo(pg *postgres.Postgres,
// 	l *logging.ZapLogger, pool *pgxpool.Pool) *GopherMartRepo {
// 	return &GopherMartRepo{
// 		pg:     pg,
// 		Logger: l,
// 		pool:   pool,
// 	}
// }

func (g *GopherMartRepo) logAndReturnError(ctx context.Context, method string, err error) error {
	msg := fmt.Sprintf("%s - %s: %v", "GopherMartRepoName", method, err)
	g.Logger.ErrorCtx(ctx, msg, zap.Error(err))
	return err
}
