package http

import (
	"encoding/json"
	"net/http"
	"strings"
)

// sendJSON отправляет JSON-ответ с указанным статус-кодом
func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendError отправляет JSON-ответ с ошибкой
func sendError(w http.ResponseWriter, message string, statusCode int) {
	sendJSON(w, map[string]interface{}{
		"error": message,
	}, statusCode)
}

// extractToken извлекает JWT токен из заголовка Authorization
// Ожидается формат: "Bearer <token>"
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Ожидается формат "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
