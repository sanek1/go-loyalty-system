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
		"token":  tokenID,
		"exp":    time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.EncryptionKey))
	if err != nil {
		return "", err
	}

	err = j.PersistToken(user, tokenID)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j TokenModel) PersistToken(user *entity.User, tokenID uuid.UUID) error {
	ctx := context.Background()
	tokenStruct := entity.Token{
		CreationDate: time.Now(),
		ID:           tokenID,
		UserID:       user.ID,
	}
	err := j.u.CreateToken(ctx, &tokenStruct)
	return err
}
