package security

import (
	"context"
	"encoding/json"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"time"
)

type TokenModel struct {
	EncryptionKey string `yaml:"jwt"`
	u             usecase.UserUseCase
	redis         *redis.Client
}

func NewJwtToken(key string, u usecase.UserUseCase, redis *redis.Client) *TokenModel {
	return &TokenModel{
		EncryptionKey: key,
		u:             u,
		redis:         redis,
	}
}

//	type Claims struct {
//	    UserID uint   `json:"user_id"`
//	    jwt.StandardClaims
//	}
const (
	sessionTimeRedis = 5
)

func (j TokenModel) GenerateToken(user *entity.User) (string, error) {
	tokenID := uuid.New()
	claims := jwt.MapClaims{
		"sub":    user.Email,
		"login":  user.Login,
		"access": user.Access,
		"id":     strconv.FormatUint(uint64(user.ID), 10),
		"token":  tokenID,                          // PK utilized to query table 'tokens'
		"exp":    time.Now().Add(time.Hour).Unix(), // Token expires in 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.EncryptionKey))
	if err != nil {
		return "", err
	}

	err = j.persistToken(user, tokenID)
	if err != nil {
		return "", err
	}

	// Сохранение сессии в Redis
	ctx := context.Background()
	data, _ := json.Marshal(user)
	j.redis.Set(ctx, strconv.FormatUint(uint64(user.ID), 10), data, sessionTimeRedis*time.Minute)

	return tokenString, nil
}

func (j TokenModel) persistToken(user *entity.User, tokenID uuid.UUID) error {
	ctx := context.Background()
	tokenStruct := entity.Token{
		CreationDate: time.Now(),
		ID:           tokenID,
		UserID:       user.ID,
	}

	err := j.u.CreateToken(ctx, &tokenStruct)
	return err
}
