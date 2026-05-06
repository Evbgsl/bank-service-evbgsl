package services

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRegisterData = errors.New("invalid register data")
	ErrUserAlreadyExists   = errors.New("user with this email or username already exists")
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
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
