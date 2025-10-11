package ctx

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := ctx.Value("userID").(string)
	if !ok {
		return uuid.UUID{}, errors.New("userID not found in context")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, errors.New("userID is not a valid uuid")
	}

	return userID, nil
}

func GetToken(ctx context.Context) string {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return ""
	}
	return token
}
