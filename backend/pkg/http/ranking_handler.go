package http

import (
	"encoding/json"
	"net/http"

	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"github.com/xartyvv/eventproject/backend/pkg/service"
)

// RankingHandler обрабатывает HTTP-запросы для ранжирования
type RankingHandler struct {
	rankingService service.RankingService
	authService    service.AuthService
}

// NewRankingHandler создает новый экземпляр RankingHandler
func NewRankingHandler(rankingService service.RankingService, authService service.AuthService) *RankingHandler {
	return &RankingHandler{
		rankingService: rankingService,
		authService:    authService,
	}
}

// SaveMatrixRequest — запрос на сохранение матрицы парных сравнений
type SaveMatrixRequest struct {
	Matrix [10][10]float64 `json:"matrix"`
}

// SaveMatrix сохраняет матрицу парных сравнений и вычисляет веса
// POST /api/ranking/matrix
func (h *RankingHandler) SaveMatrix(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString := extractToken(r)
	if tokenString == "" {
		sendError(w, "authorization required", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		sendError(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var req SaveMatrixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация матрицы
	if !isValidMatrix(req.Matrix) {
		sendError(w, "invalid matrix: diagonal must be 1, matrix[i][j] = 1/matrix[j][i]", http.StatusBadRequest)
		return
	}

	profile, err := h.rankingService.SaveRankingProfile(user.ID, req.Matrix)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Десериализуем веса для ответа
	var weights []float64
	json.Unmarshal([]byte(profile.Weights), &weights)

	sendJSON(w, map[string]interface{}{
		"message":  "ranking profile saved successfully",
		"weights":  weights,
		"criteria": domain.CriterionNames,
	}, http.StatusOK)
}

// GetWeights возвращает вычисленные веса критериев
// GET /api/ranking/weights
func (h *RankingHandler) GetWeights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString := extractToken(r)
	if tokenString == "" {
		sendError(w, "authorization required", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		sendError(w, "invalid token", http.StatusUnauthorized)
		return
	}

	profile, err := h.rankingService.GetRankingProfile(user.ID)
	if err != nil {
		sendError(w, "ranking profile not found, save your matrix first", http.StatusNotFound)
		return
	}

	var weights []float64
	var matrix [10][10]float64
	json.Unmarshal([]byte(profile.Weights), &weights)
	json.Unmarshal([]byte(profile.PairwiseMatrix), &matrix)

	sendJSON(w, map[string]interface{}{
		"weights":  weights,
		"matrix":   matrix,
		"criteria": domain.CriterionNames,
	}, http.StatusOK)
}

// GetRankedEvents возвращает ранжированный список мероприятий
// GET /api/ranking/events?category=...&start_date=...&end_date=...
func (h *RankingHandler) GetRankedEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString := extractToken(r)
	if tokenString == "" {
		sendError(w, "authorization required", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		sendError(w, "invalid token", http.StatusUnauthorized)
		return
	}

	category := r.URL.Query().Get("category")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	rankedEvents, err := h.rankingService.GetRankedEvents(user.ID, category, startDate, endDate)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"ranked_events": rankedEvents,
		"total":         len(rankedEvents),
	}, http.StatusOK)
}

// isValidMatrix проверяет корректность матрицы парных сравнений
func isValidMatrix(matrix [10][10]float64) bool {
	for i := 0; i < 10; i++ {
		// Диагональ должна быть равна 1
		if matrix[i][i] != 1.0 {
			return false
		}
		for j := 0; j < 10; j++ {
			if i == j {
				continue
			}
			// matrix[i][j] = 1/matrix[j][i]
			if matrix[i][j] <= 0 || matrix[j][i] <= 0 {
				return false
			}
			expected := 1.0 / matrix[i][j]
			relativeError := abs(expected-matrix[j][i]) / expected
			if relativeError > 0.02 { // 2% относительная ошибка
				return false
			}
		}
	}
	return true
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
