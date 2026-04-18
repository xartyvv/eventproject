package http

import (
	"encoding/json"
	"net/http"

	"github.com/xartyvv/eventproject/backend/pkg/service"
)

// AuthHandler обрабатывает HTTP-запросы для аутентификации
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRequest — структура запроса на регистрацию
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest — структура запроса на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register обрабатывает регистрацию пользователя
// POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Email == "" || req.Username == "" || req.Password == "" {
		sendError(w, "email, username, and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		sendError(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req.Email, req.Username, req.Password)
	if err != nil {
		sendError(w, err.Error(), http.StatusConflict)
		return
	}

	sendJSON(w, map[string]interface{}{
		"message": "user registered successfully",
		"user":    user,
	}, http.StatusCreated)
}

// Login обрабатывает вход пользователя
// POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		sendError(w, "email and password are required", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sendJSON(w, map[string]interface{}{
		"message": "login successful",
		"token":   token,
	}, http.StatusOK)
}

// GetMe возвращает текущего пользователя по JWT токену
// GET /api/auth/me
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString := extractToken(r)
	if tokenString == "" {
		sendError(w, "authorization header required", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		sendError(w, "invalid token", http.StatusUnauthorized)
		return
	}

	sendJSON(w, map[string]interface{}{
		"user": user,
	}, http.StatusOK)
}
