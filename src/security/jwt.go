package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/shaileshhb/equisplit/src/models"
)

var JWT_KEY = "this is a sample key, should change in prod"

func GenerateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
	})
	return token.SignedString([]byte(JWT_KEY))
}

func ValidateJWT(t string) (*models.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	sub := lo.Must(claims.GetSubject())
	exp := lo.Must(claims.GetExpirationTime())
	if exp.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}

	userId, err := uuid.Parse(sub)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Base: models.Base{
			Id: userId,
		},
	}, nil
}
