package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xartyvv/eventproject/backend/internal/database"
	httphandler "github.com/xartyvv/eventproject/backend/pkg/http"
	"github.com/xartyvv/eventproject/backend/pkg/repository"
	"github.com/xartyvv/eventproject/backend/pkg/service"
)

// Server представляет HTTP сервер приложения
type Server struct {
	router *http.ServeMux
	port   string
}

// NewServer создает новый экземпляр Server со всеми зависимостями
func NewServer() *Server {
	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(database.DB)
	eventRepo := repository.NewEventRepository(database.DB)
	rankingRepo := repository.NewRankingRepository(database.DB)
	favoriteRepo := repository.NewFavoriteRepository(database.DB)

	// Инициализируем сервисы
	authService := service.NewAuthService(userRepo)
	eventService := service.NewEventService(eventRepo)
	rankingService := service.NewRankingService(rankingRepo, eventRepo, favoriteRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo, eventRepo)

	// Инициализируем обработчики
	authHandler := httphandler.NewAuthHandler(authService)
	eventHandler := httphandler.NewEventHandler(eventService, authService)
	rankingHandler := httphandler.NewRankingHandler(rankingService, authService)
	favoriteHandler := httphandler.NewFavoriteHandler(favoriteService, authService)

	// Создаём маршрутизатор
	router := http.NewServeMux()

	// Маршруты для статических файлов (фронтенд)
	router.Handle("/", http.FileServer(http.Dir("./static")))

	// API маршрутизация
	router.HandleFunc("/api/auth/register", authHandler.Register)
	router.HandleFunc("/api/auth/login", authHandler.Login)
	router.HandleFunc("/api/auth/me", authHandler.GetMe)

	router.HandleFunc("/api/events", eventHandler.GetEvents)
	router.HandleFunc("/api/events/create", eventHandler.CreateEvent)
	router.HandleFunc("/api/events/filter", eventHandler.GetEventsByFilters)
	router.HandleFunc("/api/events/", func(w http.ResponseWriter, r *http.Request) {
		// Разделяем GET /api/events/:id и DELETE /api/events/:id
		switch r.Method {
		case http.MethodGet:
			eventHandler.GetEventByID(w, r)
		case http.MethodDelete:
			eventHandler.DeleteEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/api/ranking/matrix", rankingHandler.SaveMatrix)
	router.HandleFunc("/api/ranking/weights", rankingHandler.GetWeights)
	router.HandleFunc("/api/ranking/events", rankingHandler.GetRankedEvents)

	router.HandleFunc("/api/favorites", favoriteHandler.GetFavorites)
	router.HandleFunc("/api/favorites/add", favoriteHandler.AddFavorite)
	router.HandleFunc("/api/favorites/check/", func(w http.ResponseWriter, r *http.Request) {
		favoriteHandler.IsFavorite(w, r)
	})
	router.HandleFunc("/api/favorites/", func(w http.ResponseWriter, r *http.Request) {
		// DELETE /api/favorites/:event_id
		if r.Method == http.MethodDelete {
			favoriteHandler.RemoveFavorite(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{
		router: router,
		port:   port,
	}
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("Server starting on http://localhost%s", addr)
	return http.ListenAndServe(addr, s.router)
}
