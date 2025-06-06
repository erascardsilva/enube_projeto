package service

import (
	"errors"

	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(username, email, password string) error {
	// Check if user already exists
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return errors.New("username already exists")
	}
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Active:   true,
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(username, password string) (string, error) {
	// Find user
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
