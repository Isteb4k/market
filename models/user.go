package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Phone string `json:"phone"`
}

// GenToken - generate new token for user
func (u *User) GenToken() (string, error) {
	expiredAt := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiredAt.Unix(),
		Id:        fmt.Sprintf("%d", u.ID),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "market",
	})

	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
