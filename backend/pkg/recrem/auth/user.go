package auth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/ShareCampus/RecRem/backend/pkg/utils/logger"
	"google.golang.org/api/option"
)

const COPILOT_FIREBASE_KEY = "/etc/secret/copilot-secret"

// authenticate token via firebase
func AuthUser(ctx context.Context, token string) (*auth.Token, error) {
	opt := option.WithCredentialsFile(COPILOT_FIREBASE_KEY)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Errorf("error initializing app: %v", err)
		return nil, err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		logger.Error("error getting Auth client: %v", err)
		return nil, err
	}

	return client.VerifyIDToken(ctx, token)
}
