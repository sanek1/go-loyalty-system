package repo

import (
	"context"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/pkg/logging"
	"go-loyalty-system/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

//go:generate mockgen -source=user_pg.go -destination=./mocks/mock_user.go -package=mocks
type AuthUseCase interface {
	RegisterUser(ctx context.Context, u entity.User) error
	CreateToken(ctx context.Context, u *entity.Token) error
	GetUsers(context.Context) ([]entity.User, error)
	GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error)
	GetUserByLogin(ctx context.Context, u entity.User) (*entity.User, error)
}

func NewUserrepository(pg *postgres.Postgres, l *logging.ZapLogger, pool *pgxpool.Pool) *GopherMartRepo {
	return &GopherMartRepo{
		pg:     pg,
		Logger: l,
		pool:   pool,
	}
}
func (g *GopherMartRepo) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return g.getUser(ctx, "SELECT id, login, password, email FROM users WHERE id = $1", id)
}

func (g *GopherMartRepo) GetUserByLogin(ctx context.Context, u entity.User) (*entity.User, error) {
	return g.getUser(ctx, "SELECT id, login, password, email FROM users WHERE login = $1", u.Login)
}

func (g *GopherMartRepo) GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error) {
	return g.getUser(ctx, `SELECT id, login, password, email FROM users WHERE email = $1 and password = $2`, u.Email, u.Password)
}

func (g *GopherMartRepo) getUser(ctx context.Context, query string, args ...interface{}) (*entity.User, error) {
	row := g.pg.Pool.QueryRow(ctx, query, args...)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "Error scanning user row: %w", zap.Error(err))
		return nil, err
	}
	return user, nil
}
func (g *GopherMartRepo) RegisterUser(ctx context.Context, u entity.User) error {
	_, err := g.pg.Pool.Exec(ctx, `
		INSERT INTO users (login, email, password)
		VALUES ($1, $2, $3)
		ON CONFLICT (login) DO NOTHING
	`, u.Login, u.Email, u.Password)
	if err != nil {
		return g.logAndReturnError(ctx, "RegisterUser", err)
	}
	return nil
}

func (g *GopherMartRepo) GetUsers(ctx context.Context) ([]entity.User, error) {
	sql, _, err := g.pg.Builder.
		Select("login, email").
		From("users").
		ToSql()
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GetUsers", err)
	}

	rows, err := g.pg.Pool.Query(ctx, sql)
	if err != nil {
		return nil, g.logAndReturnError(ctx, "GetUsers", err)
	}
	defer rows.Close()
	entities := make([]entity.User, 0, _defaultEntityCap)

	for rows.Next() {
		e := entity.User{}
		if err := rows.Scan(&e.Login, &e.Email); err != nil {
			return nil, g.logAndReturnError(ctx, "GetUsers", err)
		}
		entities = append(entities, e)
	}

	return entities, nil
}
