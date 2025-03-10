package service

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/DenisEMPS/online-shop/internal/domain"
	"github.com/DenisEMPS/online-shop/internal/infastructure/kafka"
	"github.com/DenisEMPS/online-shop/internal/infastructure/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	usersTopic = "users"
)

var (
	ErrUserExists                = errors.New("user already exists")
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrInvalidToken              = errors.New("invalid token")
	ErrTokenInvalidSigningMethod = errors.New("invalid token signing method")
)

type AuthService struct {
	repo      repository.Auth
	log       *slog.Logger
	kafkaProd *kafka.Producer
}

func NewAuthService(repo repository.Auth, kafkaProd *kafka.Producer, log *slog.Logger) *AuthService {
	return &AuthService{repo: repo, kafkaProd: kafkaProd, log: log}
}

func (s *AuthService) Register(input *domain.UserCreate) (int64, error) {
	const op = "auth_service.register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", input.Email),
	)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s %w", op, err)
	}

	input.PassHash = passHash
	id, err := s.repo.Register(input)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			log.Info("user allready exists")
			return 0, ErrUserExists
		}
		log.Error("failed to register user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s %w", op, err)
	}

	log.Info("user successfuly registered")

	go func() {
		if err := s.kafkaProd.Produce(fmt.Sprintf("email %s", input.Email), usersTopic); err != nil {
			log.Error("failed to produce message about register user: %v", input.Email, slog.String("error", err.Error()))
		}
	}()
	return id, nil
}

func (s *AuthService) Login(input *domain.UserLogin) (string, error) {
	const op = "auth_service.login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", input.Email),
	)

	log.Info("attempting to login user")

	userData, err := s.repo.Login(input)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Info("user not found")
			return "", fmt.Errorf("%s %w", op, ErrUserNotFound)
		}
		log.Error("failed to get user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(userData.PassHash, []byte(input.Password)); err != nil {
		log.Info("invalid credentials", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s %w", op, ErrInvalidCredentials)
	}

	token, err := s.GenerateToken(userData)
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s %w", op, err)
	}

	log.Info("user logged in successfuly")

	return token, nil
}

func (*AuthService) GenerateToken(userData *domain.UserLoginDAO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &domain.Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
		},
		UserID:    userData.ID,
		UserEmail: userData.Email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(token string) (int64, error) {
	const op = "auth_service.parse_token"

	log := s.log.With(
		slog.String("op", op),
		slog.String("token", token),
	)

	tokenParsed, err := jwt.ParseWithClaims(token, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalidSigningMethod
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		log.Error("invalid token signing method", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrTokenInvalidSigningMethod)
	}

	if !tokenParsed.Valid {
		log.Info("token is not valid")
		return 0, ErrInvalidToken
	}

	claims, ok := tokenParsed.Claims.(*domain.Claims)
	if !ok {
		log.Warn("token claims is not in type")
		return 0, ErrInvalidToken
	}

	log.Info("token successfuly parsed", slog.String("email", claims.UserEmail))

	return claims.UserID, nil
}
