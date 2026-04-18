package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"github.com/xartyvv/eventproject/backend/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService определяет интерфейс для аутентификации
type AuthService interface {
	Register(email, username, password string) (*domain.User, error)
	Login(email, password string) (string, error)
	ValidateToken(tokenString string) (*domain.User, error)
	HashPassword(password string) (string, error)
	CheckPassword(hash, password string) bool
}

// authService — реализация AuthService
type authService struct {
	userRepo repository.UserRepository
	jwtSecret []byte
}

// NewAuthService создает новый экземпляр authService
func NewAuthService(userRepo repository.UserRepository) AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production" // fallback для разработки
	}
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(secret),
	}
}

// Register регистрирует нового пользователя
func (s *authService) Register(email, username, password string) (*domain.User, error) {
	// Проверяем, существует ли пользователь с таким email
	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Проверяем, существует ли пользователь с таким username
	if _, err := s.userRepo.GetByUsername(username); err == nil {
		return nil, errors.New("user with this username already exists")
	}

	// Хешируем пароль
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Создаём пользователя
	user := &domain.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	user.Password = "" // Не возвращаем хеш пароля
	return user, nil
}

// Login аутентифицирует пользователя и возвращает JWT токен
func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !s.CheckPassword(user.Password, password) {
		return "", errors.New("invalid email or password")
	}

	// Создаём JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Токен действует 24 часа
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken проверяет валидность JWT токена и возвращает пользователя
func (s *authService) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID := uint(claims["user_id"].(float64))
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// HashPassword хеширует пароль с помощью bcrypt
func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword проверяет соответствие пароля хешу
func (s *authService) CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
