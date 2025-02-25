package security

import (
	"context"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"time"
)

type TokenModel struct {
	EncryptionKey string `yaml:"jwt"`
	u             usecase.UserUseCase
}

func NewJwtToken(key string, u usecase.UserUseCase) *TokenModel {
	return &TokenModel{
		EncryptionKey: key,
		u:             u,
	}
}

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
