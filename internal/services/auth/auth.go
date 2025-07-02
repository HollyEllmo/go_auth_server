package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/HollyEllmo/go_auth_server/internal/domain/models"
	"github.com/HollyEllmo/go_auth_server/internal/lib/jwt"
	"github.com/HollyEllmo/go_auth_server/internal/lib/logger/sl"
	"github.com/HollyEllmo/go_auth_server/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log *slog.Logger
	userSaver UserSaver
	userProvider UserProvider
	appProvider AppProvider
	tokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID = errors.New("invalid app ID")
	ErrUserExists = errors.New("user already exists")
)

// New returns a new instance of Auth service.
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log: log.With(slog.String("service", "auth")),
		userSaver: userSaver,
		userProvider: userProvider,
		appProvider: appProvider,
		tokenTTL: tokenTTL,
	}
}

// Login authenticates a user by email and password, and returns a token if successful.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (token string, err error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: user not found: %w", op, ErrInvalidCredentials) 
		}

		a.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid password", sl.Err(err))
		return "", fmt.Errorf("%s: invalid credentials: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: failed to get app: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err = jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to create token", sl.Err(err))
		return "", fmt.Errorf("%s: failed to create token: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registrating new user",)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			log.Warn("user already exists", slog.String("email", email))
			return 0, fmt.Errorf("%s: user already exists: %w", op, ErrUserExists)
		}
		log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: failed to generate password hash: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: failed to save user: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", sl.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		log.Error("failed to check if user is admin", sl.Err(err))
		return false, fmt.Errorf("%s: failed to check if user is admin: %w", op, err)
	}

	log.Info("admin check completed", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
