package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Event struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Location    string `json:"location"`
}

type AddFavoriteRequest struct {
	EventID uint `json:"event_id"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var events = []Event{
	{ID: 1, Title: "Концерт уличной музыки", Description: "Летний концерт в парке с живыми выступлениями.", Date: "2026-05-01", Location: "Парк Горького"},
	{ID: 2, Title: "Выставка современного искусства", Description: "Новые работы молодых художников.", Date: "2026-05-10", Location: "Центр культуры"},
	{ID: 3, Title: "Онлайн мастер-класс по кулинарии", Description: "Готовим десерты вместе с шефом.", Date: "2026-05-15", Location: "Онлайн"},
}

var favorites = map[uint]bool{}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/events", handleEvents)
	mux.HandleFunc("/api/favorites", handleFavorites)
	mux.HandleFunc("/api/favorites/", handleFavoriteDelete)
	mux.HandleFunc("/api/login", handleAuthLogin)
	mux.HandleFunc("/api/register", handleAuthRegister)

	port := 8081
	fmt.Printf("Backend is running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), cors(mux)))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sendJSON(w, events)
}

func handleFavorites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sendFavorites(w)
	case http.MethodPost:
		addFavorite(w, r)
	default:
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleFavoriteDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/api/favorites/"):]
	eventID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		sendError(w, "invalid event id", http.StatusBadRequest)
		return
	}

	delete(favorites, uint(eventID))
	sendJSON(w, map[string]interface{}{"message": "removed"})
}

func handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sendJSON(w, map[string]interface{}{
		"message": "login endpoint stub: auth not implemented",
		"email":   req.Email,
	})
}

func handleAuthRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sendJSON(w, map[string]interface{}{
		"message": "register endpoint stub: auth not implemented",
		"email":   req.Email,
	})
}

func sendFavorites(w http.ResponseWriter) {
	var list []Event
	for _, event := range events {
		if favorites[event.ID] {
			list = append(list, event)
		}
	}
	sendJSON(w, map[string]interface{}{"favorites": list})
}

func addFavorite(w http.ResponseWriter, r *http.Request) {
	var req AddFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.EventID == 0 {
		sendError(w, "event_id is required", http.StatusBadRequest)
		return
	}

	if !eventExists(req.EventID) {
		sendError(w, "event not found", http.StatusNotFound)
		return
	}

	favorites[req.EventID] = true
	sendJSON(w, map[string]interface{}{"message": "added"})
}

func eventExists(id uint) bool {
	for _, event := range events {
		if event.ID == id {
			return true
		}
	}
	return false
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	sendJSON(w, ErrorResponse{Error: message})
}
