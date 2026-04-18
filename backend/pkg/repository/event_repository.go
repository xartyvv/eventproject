package repository

import (
	"time"

	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"gorm.io/gorm"
)

// EventRepository определяет интерфейс для работы с мероприятиями
type EventRepository interface {
	Create(event *domain.Event) error
	GetByID(id uint) (*domain.Event, error)
	GetAll(limit, offset int) ([]domain.Event, error)
	GetByCategory(category string, limit, offset int) ([]domain.Event, error)
	GetByDateRange(start, end time.Time, limit, offset int) ([]domain.Event, error)
	GetByFilters(category string, start, end time.Time, limit, offset int) ([]domain.Event, error)
	GetByCreatorID(creatorID uint) ([]domain.Event, error)
	Update(event *domain.Event) error
	Delete(id uint) error
	GetAllForRanking() ([]domain.Event, error)
}

// eventRepository — реализация EventRepository с использованием GORM
type eventRepository struct {
	db *gorm.DB
}

// NewEventRepository создает новый экземпляр eventRepository
func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *domain.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) GetByID(id uint) (*domain.Event, error) {
	var event domain.Event
	err := r.db.Preload("Creator").First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetAll(limit, offset int) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.Preload("Creator").
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

func (r *eventRepository) GetByCategory(category string, limit, offset int) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.Where("category = ?", category).
		Preload("Creator").
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

func (r *eventRepository) GetByDateRange(start, end time.Time, limit, offset int) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.Where("date BETWEEN ? AND ?", start, end).
		Preload("Creator").
		Order("date ASC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

func (r *eventRepository) GetByFilters(category string, start, end time.Time, limit, offset int) ([]domain.Event, error) {
	var events []domain.Event
	query := r.db.Preload("Creator").Model(&domain.Event{})

	// Фильтр по категории
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Фильтр по дате
	if !start.IsZero() && !end.IsZero() {
		query = query.Where("date BETWEEN ? AND ?", start, end)
	}

	err := query.Order("date ASC").Limit(limit).Offset(offset).Find(&events).Error
	return events, err
}

func (r *eventRepository) GetByCreatorID(creatorID uint) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.Where("creator_id = ?", creatorID).
		Order("date DESC").
		Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(event *domain.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Event{}, id).Error
}

func (r *eventRepository) GetAllForRanking() ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.Order("date ASC").Find(&events).Error
	return events, err
}
