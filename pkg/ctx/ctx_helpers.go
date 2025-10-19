package ctx

import (
	"context"
	"errors"

	"github.com/Melodia-IS2/melodia-go-utils/pkg/region"
	"github.com/google/uuid"
)

type ContextState string

const (
	ContextStateActive      ContextState = "active"
	ContextStateBlocked     ContextState = "blocked"
	ContextStateMissingInfo ContextState = "missing_info"
	ContextStateNoSession   ContextState = "no_session"
)

type ContextRol string

const (
	ContextRolAdmin  ContextRol = "admin"
	ContextRolArtist ContextRol = "artist"
	ContextRolUser   ContextRol = "user"
	ContextRolGuest  ContextRol = "guest"
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

func GetState(ctx context.Context) ContextState {
	state, ok := ctx.Value("state").(string)
	if !ok {
		return ContextStateNoSession
	}
	return ContextState(state)
}

func GetRol(ctx context.Context) string {
	rol, ok := ctx.Value("rol").(string)
	if !ok {
		return ""
	}
	return rol
}

func GetRegion(ctx context.Context) region.Region {
	reg, ok := ctx.Value("region").(region.Region)
	if !ok {
		return region.Global
	}
	return reg
}
