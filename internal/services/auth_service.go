package services

import (
	"errors"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 24 * time.Hour

var (
	ErrInvalidRegisterData = errors.New("invalid register data")
	ErrUserAlreadyExists   = errors.New("user with this email or username already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password

	if !isValidUsername(username) {
		return nil, ErrInvalidRegisterData
	}

	if !isValidEmail(email) {
		return nil, ErrInvalidRegisterData
	}

	if !isValidPassword(password) {
		return nil, ErrInvalidRegisterData
	}

	exists, err := s.userRepo.ExistsByEmailOrUsername(email, username)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}

		return nil, err
	}

	return &models.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "user registered successfully",
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password

	if !isValidEmail(email) || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:            token,
		TokenType:        "Bearer",
		ExpiresInSeconds: int64(tokenTTL.Seconds()),
	}, nil
}

func (s *AuthService) generateJWT(userID int64) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

func isValidUsername(username string) bool {
	length := utf8.RuneCountInString(username)
	return length >= 3 && length <= 50
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(password string) bool {
	return len(password) >= 6
}
