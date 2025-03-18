package auth

import (
	"Bybit_Pet_Project/internal/domain/models"
	"Bybit_Pet_Project/internal/lib/jwt"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log/slog"
	"time"
)

type UserStorage interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
	User(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int64) (models.App, error)
}

type Auth struct {
	log           *slog.Logger
	userInterface UserStorage
	appProvider   AppProvider
	tokenTTL      time.Duration
}

func New(log *slog.Logger, userInterface UserStorage, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:           log,
		userInterface: userInterface,
		appProvider:   appProvider,
		tokenTTL:      tokenTTL,
	}
}
func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {
	const op = "Auth.RegisterNewUser"
	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("faliled to generate password hash")
		return 0, fmt.Errorf("%s, %w", op, err)
	}
	id, err := a.userInterface.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("faliled to save user")
		return 0, fmt.Errorf("%s, %w", op, err)
	}
	return id, nil
}

func (a *Auth) Login(ctx context.Context, email string, pass string, appID int) (string, error) {
	const op = "Auth.Login"
	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("attemting to login user")

	user, err := a.userInterface.User(ctx, email)
	if err != nil {
		log.Error("faliled to get user")
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass)); err != nil {
		log.Error("faliled to compare password")
		return "", fmt.Errorf("%s", op)
	}
	app, err := a.appProvider.App(ctx, int64(appID))
	if err != nil {
		log.Error("faliled to get app")
	}
	log.Info("login sucsess")
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token")

		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
