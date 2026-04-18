package service

import (
	"time"

	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"github.com/xartyvv/eventproject/backend/pkg/repository"
)

// EventService определяет интерфейс для работы с мероприятиями
type EventService interface {
	Create(event *domain.Event) (*domain.Event, error)
	GetByID(id uint) (*domain.Event, error)
	GetAll(page, pageSize int) ([]domain.Event, int, error)
	GetByFilters(category string, start, end time.Time, page, pageSize int) ([]domain.Event, int, error)
	GetByCreatorID(creatorID uint) ([]domain.Event, error)
	Update(event *domain.Event) (*domain.Event, error)
	Delete(id uint, userID uint) error
	GetAllForRanking() ([]domain.Event, error)
}

// eventService — реализация EventService
type eventService struct {
	eventRepo repository.EventRepository
}

// NewEventService создает новый экземпляр eventService
func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) Create(event *domain.Event) (*domain.Event, error) {
	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) GetByID(id uint) (*domain.Event, error) {
	return s.eventRepo.GetByID(id)
}

func (s *eventService) GetAll(page, pageSize int) ([]domain.Event, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	events, err := s.eventRepo.GetAll(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Подсчитываем общее количество
	allEvents, _ := s.eventRepo.GetAll(100000, 0)
	total := len(allEvents)

	return events, total, nil
}

func (s *eventService) GetByFilters(category string, start, end time.Time, page, pageSize int) ([]domain.Event, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	events, err := s.eventRepo.GetByFilters(category, start, end, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Подсчитываем общее количество с фильтрами
	allFiltered, _ := s.eventRepo.GetByFilters(category, start, end, 100000, 0)
	total := len(allFiltered)

	return events, total, nil
}

func (s *eventService) GetByCreatorID(creatorID uint) ([]domain.Event, error) {
	return s.eventRepo.GetByCreatorID(creatorID)
}

func (s *eventService) Update(event *domain.Event) (*domain.Event, error) {
	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) Delete(id uint, userID uint) error {
	// Проверяем, что пользователь является создателем мероприятия
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return err
	}
	if event.CreatorID != userID {
		return &AuthorizationError{Message: "only the creator can delete this event"}
	}
	return s.eventRepo.Delete(id)
}

func (s *eventService) GetAllForRanking() ([]domain.Event, error) {
	return s.eventRepo.GetAllForRanking()
}

// AuthorizationError — ошибка авторизации
type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}
