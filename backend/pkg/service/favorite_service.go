package service

import (
	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"github.com/xartyvv/eventproject/backend/pkg/repository"
)

// FavoriteService определяет интерфейс для работы с избранным
type FavoriteService interface {
	AddFavorite(userID, eventID uint) error
	RemoveFavorite(userID, eventID uint) error
	GetFavorites(userID uint, page, pageSize int) ([]domain.Favorite, int, error)
	IsFavorite(userID, eventID uint) bool
}

// favoriteService — реализация FavoriteService
type favoriteService struct {
	favoriteRepo repository.FavoriteRepository
	eventRepo    repository.EventRepository
}

// NewFavoriteService создает новый экземпляр favoriteService
func NewFavoriteService(favoriteRepo repository.FavoriteRepository, eventRepo repository.EventRepository) FavoriteService {
	return &favoriteService{
		favoriteRepo: favoriteRepo,
		eventRepo:    eventRepo,
	}
}

func (s *favoriteService) AddFavorite(userID, eventID uint) error {
	// Проверяем, существует ли мероприятие
	_, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return &NotFoundError{Message: "event not found"}
	}

	// Проверяем, не добавлено ли уже в избранное
	if s.favoriteRepo.IsFavorite(userID, eventID) {
		return &AlreadyExistsError{Message: "event is already in favorites"}
	}

	favorite := &domain.Favorite{
		UserID:  userID,
		EventID: eventID,
	}

	return s.favoriteRepo.Add(favorite)
}

func (s *favoriteService) RemoveFavorite(userID, eventID uint) error {
	return s.favoriteRepo.Remove(userID, eventID)
}

func (s *favoriteService) GetFavorites(userID uint, page, pageSize int) ([]domain.Favorite, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	favorites, err := s.favoriteRepo.GetByUserID(userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Подсчитываем общее количество
	allFavorites, _ := s.favoriteRepo.GetByUserID(userID, 100000, 0)
	total := len(allFavorites)

	return favorites, total, nil
}

func (s *favoriteService) IsFavorite(userID, eventID uint) bool {
	return s.favoriteRepo.IsFavorite(userID, eventID)
}

// NotFoundError — ошибка "не найдено"
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// AlreadyExistsError — ошибка "уже существует"
type AlreadyExistsError struct {
	Message string
}

func (e *AlreadyExistsError) Error() string {
	return e.Message
}
