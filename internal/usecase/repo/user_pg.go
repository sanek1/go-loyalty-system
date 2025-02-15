package repo

import (
	"context"
	"go-loyalty-system/internal/entity"

	"go.uber.org/zap"
)

func (g *GopherMartRepo) GetUserByEmail(ctx context.Context, u entity.User) (*entity.User, error) {
	const query = `SELECT id, login, password, email FROM users WHERE email = $1 and password = $2`
	row := g.pg.Pool.QueryRow(ctx, query, u.Email, u.Password)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "Error scanning user row: %w", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (g *GopherMartRepo) GetUserByID(ctx context.Context, u entity.User) (*entity.User, error) {
	const query = `SELECT id, login, password, email FROM users WHERE id = $1 `
	row := g.pg.Pool.QueryRow(ctx, query, u.ID)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "Error scanning user row: %w", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (g *GopherMartRepo) GetUserByLogin(ctx context.Context, u entity.User) (*entity.User, error) {
	const query = `SELECT id, login, password, email FROM users WHERE login = $1 `
	row := g.pg.Pool.QueryRow(ctx, query, u.Login)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
	if err != nil {
		g.Logger.ErrorCtx(ctx, "Error scanning user row: %w", zap.Error(err))
		return nil, err
	}
	return user, nil
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

func (g *GopherMartRepo) RegisterUser(ctx context.Context, u entity.User) error {
	sql, args, err := g.pg.Builder.
		Insert("users").
		Columns("login, email", "password").
		Values(u.Login, u.Email, u.Password).
		ToSql()
	if err != nil {
		return g.logAndReturnError(ctx, "RegisterUser", err)
	}

	_, err = g.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return g.logAndReturnError(ctx, "RegisterUser", err)
	}
	_ = g.SetBalance(ctx, u)
	return nil
}

func (g *GopherMartRepo) SetBalance(ctx context.Context, u entity.User) error {
	user, _ := g.GetUserByLogin(ctx, u)
	if user == nil {
		return nil
	}
	sql, args, err := g.pg.Builder.
		Insert("balance").
		Columns("user_id, current_balance, withdrawn").
		Values(user.ID, 10000, 0).
		ToSql()
	if err != nil {
		return g.logAndReturnError(ctx, "SetBalance", err)
	}

	_, err = g.pg.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return g.logAndReturnError(ctx, "SetBalance", err)
	}

	return nil

}
