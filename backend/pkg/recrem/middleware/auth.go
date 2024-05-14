package middleware

import (
	"context"
	"fmt"
	"net/http"

	fbAuth "firebase.google.com/go/v4/auth"
	"github.com/ShareCampus/RecRem/backend/pkg/recrem/auth"
	"github.com/ShareCampus/RecRem/backend/pkg/utils/httpjson"
	"github.com/ShareCampus/RecRem/backend/pkg/utils/logger"
)

type contextKey string

const (
	AUTH                            = "Authorization"
	CONTEXT_USERINFO_KEY contextKey = "UserInfo"
)

// Auth is a middleware to verify the token in the request header
// and obtain the user information from the token
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.Header.Values(AUTH)
		if len(values) == 0 {
			logger.Debugf("no authorization header found")
			httpjson.ReturnError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Verify the token and obtain its user information
		userInfo, err := auth.AuthUser(r.Context(), values[0])
		if err != nil {
			logger.Debugf("failed to authenticate user: %s", err.Error())
			httpjson.ReturnError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		ctxWithUserInfo := context.WithValue(r.Context(), CONTEXT_USERINFO_KEY, userInfo)
		r = r.WithContext(ctxWithUserInfo)

		next.ServeHTTP(w, r)
	})
}

// GetTokenFromContext returns the user information from the context
// This function is available when the Auth middleware is used
func GetTokenFromContext(ctx context.Context) (*fbAuth.Token, error) {
	userInfo, exists := ctx.Value(CONTEXT_USERINFO_KEY).(*fbAuth.Token)
	if !exists {
		return nil, fmt.Errorf("failed to get user info from context")
	}

	return userInfo, nil
}
