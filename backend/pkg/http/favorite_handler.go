package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/xartyvv/eventproject/backend/pkg/service"
)

// FavoriteHandler обрабатывает HTTP-запросы для избранного
type FavoriteHandler struct {
	favoriteService service.FavoriteService
	authService     service.AuthService
}

// NewFavoriteHandler создает новый экземпляр FavoriteHandler
func NewFavoriteHandler(favoriteService service.FavoriteService, authService service.AuthService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
		authService:     authService,
	}
}

// AddFavoriteRequest — запрос на добавление в избранное
type AddFavoriteRequest struct {
	EventID uint `json:"event_id"`
}

// AddFavorite добавляет мероприятие в избранное
// POST /api/favorites
func (h *FavoriteHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
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

	var req AddFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.EventID == 0 {
		sendError(w, "event_id is required", http.StatusBadRequest)
		return
	}

	if err := h.favoriteService.AddFavorite(user.ID, req.EventID); err != nil {
		switch err.(type) {
		case *service.NotFoundError:
			sendError(w, err.Error(), http.StatusNotFound)
		case *service.AlreadyExistsError:
			sendError(w, err.Error(), http.StatusConflict)
		default:
			sendError(w, "failed to add to favorites", http.StatusInternalServerError)
		}
		return
	}

	sendJSON(w, map[string]interface{}{
		"message":  "event added to favorites",
		"event_id": req.EventID,
	}, http.StatusCreated)
}

// RemoveFavorite удаляет мероприятие из избранного
// DELETE /api/favorites/:event_id
func (h *FavoriteHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	idStr := r.URL.Path[len("/api/favorites/"):]
	eventID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		sendError(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	if err := h.favoriteService.RemoveFavorite(user.ID, uint(eventID)); err != nil {
		sendError(w, "failed to remove from favorites", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"message":  "event removed from favorites",
		"event_id": eventID,
	}, http.StatusOK)
}

// GetFavorites возвращает список избранных мероприятий
// GET /api/favorites?page=1&page_size=20
func (h *FavoriteHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
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

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	favorites, total, err := h.favoriteService.GetFavorites(user.ID, page, pageSize)
	if err != nil {
		sendError(w, "failed to get favorites", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"favorites": favorites,
		"total":     total,
		"page":      page,
	}, http.StatusOK)
}

// IsFavorite проверяет, находится ли мероприятие в избранном
// GET /api/favorites/check/:event_id
func (h *FavoriteHandler) IsFavorite(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.URL.Path[len("/api/favorites/check/"):]
	eventID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		sendError(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	isFav := h.favoriteService.IsFavorite(user.ID, uint(eventID))

	sendJSON(w, map[string]interface{}{
		"is_favorite": isFav,
		"event_id":    eventID,
	}, http.StatusOK)
}
