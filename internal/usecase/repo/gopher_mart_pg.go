package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-loyalty-system/internal/entity"
	"go-loyalty-system/pkg/logging"
	"go-loyalty-system/pkg/postgres"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const _defaultEntityCap = 64
const GopherMartRepoName = "GopherMartRepo"

type GopherMartRepo struct {
	pg     *postgres.Postgres
	Logger *logging.ZapLogger
	redis  *redis.Client
}

func NewUserRepo(pg *postgres.Postgres, redis *redis.Client, l *logging.ZapLogger) *GopherMartRepo {
	return &GopherMartRepo{
		pg:     pg,
		Logger: l,
		redis:  redis,
	}
}

func (g *GopherMartRepo) logAndReturnError(ctx context.Context, method string, err error) error {
	msg := fmt.Sprintf("%s - %s: %v", "GopherMartRepoName", method, err)
	g.Logger.ErrorCtx(ctx, msg, zap.Error(err))
	return err
}

func (g *GopherMartRepo) GetCurrentUser(id uint) (user *entity.User, err error) {
	ctx := context.Background()
	var item []byte
	if item, err = g.redis.Get(ctx, fmt.Sprint(id)).Bytes(); err != nil {
		if err == redis.Nil {
			return nil, http.ErrNoCookie
		}
		return nil, g.logAndReturnError(ctx, "GetCurrentUser - redis", err)
	}
	if err = json.Unmarshal(item, &user); err != nil {
		return nil, g.logAndReturnError(ctx, "GetCurrentUser - json.Unmarshal", err)
	}
	return user, nil
}
