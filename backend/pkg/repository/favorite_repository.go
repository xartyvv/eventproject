package repository

import (
	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"gorm.io/gorm"
)

// FavoriteRepository определяет интерфейс для работы с избранным
type FavoriteRepository interface {
	Add(favorite *domain.Favorite) error
	Remove(userID, eventID uint) error
	GetByUserID(userID uint, limit, offset int) ([]domain.Favorite, error)
	IsFavorite(userID, eventID uint) bool
}

// favoriteRepository — реализация FavoriteRepository с использованием GORM
type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository создает новый экземпляр favoriteRepository
func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Add(favorite *domain.Favorite) error {
	return r.db.Create(favorite).Error
}

func (r *favoriteRepository) Remove(userID, eventID uint) error {
	return r.db.Where("user_id = ? AND event_id = ?", userID, eventID).Delete(&domain.Favorite{}).Error
}

func (r *favoriteRepository) GetByUserID(userID uint, limit, offset int) ([]domain.Favorite, error) {
	var favorites []domain.Favorite
	err := r.db.Preload("Event").Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	return favorites, nil
}

func (r *favoriteRepository) IsFavorite(userID, eventID uint) bool {
	var count int64
	r.db.Model(&domain.Favorite{}).Where("user_id = ? AND event_id = ?", userID, eventID).Count(&count)
	return count > 0
}
