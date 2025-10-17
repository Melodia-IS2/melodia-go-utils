package auth

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Melodia-IS2/melodia-go-utils/pkg/ctx"
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
		token, claims, err := ReadClaims(r, a.JWTSecretKey)
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
		ctx = context.WithValue(ctx, "token", token)

		return next(w, r.WithContext(ctx))
	}
}

func (a *AuthMiddleware) CheckKeyValue(c context.Context, key string, validValues []string) bool {
	value := c.Value(key).(string)

	return slices.Contains(validValues, strings.ToUpper(value))
}

type MiddlewareValidation func(r *http.Request) error

type Builder struct {
	auth   *AuthMiddleware
	checks []MiddlewareValidation
}

func (a *AuthMiddleware) NewBuilder() *Builder {
	return &Builder{auth: a, checks: []MiddlewareValidation{}}
}

func (b *Builder) WithState(allowed ...ctx.ContextState) *Builder {
	upper := make([]string, len(allowed))
	for i, v := range allowed {
		upper[i] = strings.ToUpper(string(v))
	}
	b.checks = append(b.checks, func(r *http.Request) error {
		stateVal := ctx.GetState(r.Context())
		if !slices.Contains(upper, strings.ToUpper(string(stateVal))) {
			return errors.NewUnauthorizedError("state not allowed")
		}
		return nil
	})
	return b
}

func (b *Builder) WithRol(allowed ...ctx.ContextRol) *Builder {
	upper := make([]string, len(allowed))
	for i, v := range allowed {
		upper[i] = strings.ToUpper(string(v))
	}
	b.checks = append(b.checks, func(r *http.Request) error {
		roleVal := ctx.GetRol(r.Context())
		if !slices.Contains(upper, strings.ToUpper(string(roleVal))) {
			return errors.NewUnauthorizedError("role not allowed")
		}
		return nil
	})
	return b
}

func (b *Builder) WithClaim(claimKey string, predicate func(value any) bool, errMsg string) *Builder {
	b.checks = append(b.checks, func(r *http.Request) error {
		value := r.Context().Value(claimKey)
		if !predicate(value) {
			return errors.NewUnauthorizedError(errMsg)
		}
		return nil
	})
	return b
}

func (b *Builder) WithCustom(fn func(r *http.Request) error) *Builder {
	b.checks = append(b.checks, fn)
	return b
}

func (b *Builder) Build(next func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) error {
	return b.auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) error {
		for _, check := range b.checks {
			if err := check(r); err != nil {
				return err
			}
		}
		return next(w, r)
	})
}
