package auth

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"market/models"
	"market/repositories"
	"net/http"
	"os"
	"strconv"
)

var NoUserInContextErr = errors.New("no user in context")

const CurrentUserKey = "currentUser"

// CurrentUserMiddleware - middleware to extract Authorization header from request and match current user
func CurrentUserMiddleware(usersRepo repositories.Users) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, ok := extractClaims(token)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			userID, err := strconv.Atoi(claims["jti"].(string))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := usersRepo.GetUserByID(userID)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), CurrentUserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := os.Getenv("JWT_SECRET")
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

// GetCurrentUserFromCTX - get current user from context
func GetCurrentUserFromCTX(ctx context.Context) (*models.User, error) {
	if ctx.Value(CurrentUserKey) == nil {
		return nil, NoUserInContextErr
	}

	user, ok := ctx.Value(CurrentUserKey).(models.User)
	if !ok || user.ID == 0 {
		return nil, NoUserInContextErr
	}

	return &user, nil
}
