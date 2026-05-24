package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"github.com/xartyvv/eventproject/backend/pkg/service"
)

// EventHandler обрабатывает HTTP-запросы для мероприятий
type EventHandler struct {
	eventService service.EventService
	authService  service.AuthService
}

// NewEventHandler создает новый экземпляр EventHandler
func NewEventHandler(eventService service.EventService, authService service.AuthService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		authService:  authService,
	}
}

// CreateEventRequest — структура запроса на создание мероприятия
type CreateEventRequest struct {
	Title                string  `json:"title"`
	Description          string  `json:"description"`
	Date                 string  `json:"date"`
	Location             string  `json:"location"`
	Category             string  `json:"category"`
	Cost                 float64 `json:"cost"`
	Duration             float64 `json:"duration"`
	Capacity             float64 `json:"capacity"`
	IsWeekend            bool    `json:"is_weekend"`
	IsOnline             bool    `json:"is_online"`
	AgeRestriction       int     `json:"age_restriction"`
	RequiresRegistration bool    `json:"requires_registration"`
	OrganizerRating      float64 `json:"organizer_rating"`
	TimeOfDay            int     `json:"time_of_day"`
	Interactivity        float64 `json:"interactivity"`
}

// CreateEvent создает новое мероприятие
// POST /api/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Аутентификация
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

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Парсим дату
	eventDate, err := time.Parse("2006-01-02T15:04", req.Date)
	if err != nil {
		eventDate, err = time.Parse("2006-01-02 15:04", req.Date)
		if err != nil {
			eventDate, err = time.Parse("2006-01-02", req.Date)
			if err != nil {
				sendError(w, "invalid date format, use YYYY-MM-DD, YYYY-MM-DD HH:MM or YYYY-MM-DDTHH:MM", http.StatusBadRequest)
				return
			}
		}
	}

	event := &domain.Event{
		Title:                req.Title,
		Description:          req.Description,
		Date:                 eventDate,
		Location:             req.Location,
		Category:             req.Category,
		CreatorID:            user.ID,
		Cost:                 req.Cost,
		Duration:             req.Duration,
		Capacity:             req.Capacity,
		IsWeekend:            req.IsWeekend,
		IsOnline:             req.IsOnline,
		AgeRestriction:       req.AgeRestriction,
		RequiresRegistration: req.RequiresRegistration,
		OrganizerRating:      req.OrganizerRating,
		TimeOfDay:            req.TimeOfDay,
		Interactivity:        req.Interactivity,
	}

	event, err = h.eventService.Create(event)
	if err != nil {
		sendError(w, "failed to create event", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"message": "event created successfully",
		"event":   event,
	}, http.StatusCreated)
}

// GetEvents возвращает список всех мероприятий с пагинацией
// GET /api/events?page=1&page_size=20
func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	events, total, err := h.eventService.GetAll(page, pageSize)
	if err != nil {
		sendError(w, "failed to get events", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"events": events,
		"total":  total,
		"page":   page,
	}, http.StatusOK)
}

// GetEventByID возвращает мероприятие по ID
// GET /api/events/:id
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/events/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		sendError(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := h.eventService.GetByID(uint(id))
	if err != nil {
		sendError(w, "event not found", http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]interface{}{
		"event": event,
	}, http.StatusOK)
}

// GetEventsByFilters возвращает мероприятия с фильтрами
// GET /api/events/filter?category=...&start_date=...&end_date=...&page=1&page_size=20
func (h *EventHandler) GetEventsByFilters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
	}

	events, total, err := h.eventService.GetByFilters(category, startDate, endDate, page, pageSize)
	if err != nil {
		sendError(w, "failed to get events", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"events": events,
		"total":  total,
		"page":   page,
	}, http.StatusOK)
}

// GetMyEvents возвращает мероприятия текущего пользователя
// GET /api/events/mine
func (h *EventHandler) GetMyEvents(w http.ResponseWriter, r *http.Request) {
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

	events, err := h.eventService.GetByCreatorID(user.ID)
	if err != nil {
		sendError(w, "failed to get events", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"events": events,
		"total":  len(events),
	}, http.StatusOK)
}

// UpdateEvent обновляет мероприятие (только создатель)
// PUT /api/events/:id
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	idStr := r.URL.Path[len("/api/events/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		sendError(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	existingEvent, err := h.eventService.GetByID(uint(id))
	if err != nil {
		sendError(w, "event not found", http.StatusNotFound)
		return
	}

	if existingEvent.CreatorID != user.ID {
		sendError(w, "only the creator can update this event", http.StatusForbidden)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	eventDate, err := parseDate(req.Date)
	if err != nil {
		sendError(w, "invalid date format, use YYYY-MM-DD, YYYY-MM-DD HH:MM or YYYY-MM-DDTHH:MM", http.StatusBadRequest)
		return
	}

	existingEvent.Title = req.Title
	existingEvent.Description = req.Description
	existingEvent.Date = eventDate
	existingEvent.Location = req.Location
	existingEvent.Category = req.Category
	existingEvent.Cost = req.Cost
	existingEvent.Duration = req.Duration
	existingEvent.Capacity = req.Capacity
	existingEvent.IsWeekend = req.IsWeekend
	existingEvent.IsOnline = req.IsOnline
	existingEvent.AgeRestriction = req.AgeRestriction
	existingEvent.RequiresRegistration = req.RequiresRegistration
	existingEvent.OrganizerRating = req.OrganizerRating
	existingEvent.TimeOfDay = req.TimeOfDay
	existingEvent.Interactivity = req.Interactivity

	updatedEvent, err := h.eventService.Update(existingEvent)
	if err != nil {
		sendError(w, "failed to update event", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"message": "event updated successfully",
		"event":   updatedEvent,
	}, http.StatusOK)
}

func parseDate(dateStr string) (time.Time, error) {
	eventDate, err := time.Parse("2006-01-02T15:04", dateStr)
	if err != nil {
		eventDate, err = time.Parse("2006-01-02 15:04", dateStr)
		if err != nil {
			eventDate, err = time.Parse("2006-01-02", dateStr)
		}
	}
	return eventDate, err
}

// DeleteEvent удаляет мероприятие (только создатель)
// DELETE /api/events/:id
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.URL.Path[len("/api/events/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		sendError(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	if err := h.eventService.Delete(uint(id), user.ID); err != nil {
		if _, ok := err.(*service.AuthorizationError); ok {
			sendError(w, err.Error(), http.StatusForbidden)
			return
		}
		sendError(w, "event not found", http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]interface{}{
		"message": "event deleted successfully",
	}, http.StatusOK)
}
