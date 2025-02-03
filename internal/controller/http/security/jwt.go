package security

import (
	"context"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"

	"github.com/golang-jwt/jwt/v5"

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

	// Generate a new v4 UUID as identifier for token. This value will be used as primary key
	//tokenID :=

	// Define the JWT claims
	claims := jwt.MapClaims{
		"sub":    user.Email,
		"access": user.Access,
		//"id":     tokenID,                          // PK utilized to query table 'tokens'
		"exp": time.Now().Add(time.Hour).Unix(), // Token expires in 1 hour
	}

	// CreateToken the JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our encryption key derived from JWT configuration
	tokenString, err := token.SignedString([]byte(j.EncryptionKey))
	if err != nil {
		return "", err
	}

	// Insert struct containing necessary metadata in order to query tokens based from claims
	err = j.persistToken(user)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j TokenModel) persistToken(user *entity.User) error {
	ctx := context.Background()
	tokenStruct := entity.Token{
		CreationDate: time.Now(),
		UserID:       user.ID,
	}

	err := j.u.CreateToken(ctx, &tokenStruct)
	return err
}
