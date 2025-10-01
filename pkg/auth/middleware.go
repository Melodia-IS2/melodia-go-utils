package auth

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Melodia-IS2/melodia-go-utils/pkg/errors"
)

type AuthMiddleware struct {
	JWTSecretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{JWTSecretKey: secretKey}
}

func (a *AuthMiddleware) AuthMiddleware(next func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		claims, err := ReadClaims(r, a.JWTSecretKey)
		if err != nil {
			return errors.NewUnauthorizedError(err.Error())
		}

		if claims == nil {
			return errors.NewUnauthorizedError("invalid token")
		}
		expirationDateStr, ok := claims["expiration_date"].(string)
		if !ok {
			return errors.NewUnauthorizedError("invalid token")
		}

		expirationDate, err := time.Parse(time.RFC3339, expirationDateStr)
		if err != nil {
			return errors.NewUnauthorizedError("invalid token")
		}

		if expirationDate.Before(time.Now()) {
			return errors.NewUnauthorizedError("expired token")
		}

		ctx := context.WithValue(r.Context(), "userID", claims["user_id"])
		ctx = context.WithValue(ctx, "sessionID", claims["session_id"])
		ctx = context.WithValue(ctx, "expirationDate", expirationDate)
		ctx = context.WithValue(ctx, "state", strings.ToUpper(claims["state"].(string)))
		ctx = context.WithValue(ctx, "rol", strings.ToUpper(claims["rol"].(string)))

		return next(w, r.WithContext(ctx))
	}
}

func (a *AuthMiddleware) CheckKeyValue(c context.Context, key string, validValues []string) bool {
	value := c.Value(key).(string)

	return slices.Contains(validValues, strings.ToUpper(value))
}
