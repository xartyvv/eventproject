package repository

import (
	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"gorm.io/gorm"
)

// RankingRepository определяет интерфейс для работы с профилями ранжирования
type RankingRepository interface {
	CreateOrUpdate(profile *domain.RankingProfile) error
	GetByUserID(userID uint) (*domain.RankingProfile, error)
}

// rankingRepository — реализация RankingRepository
type rankingRepository struct {
	db *gorm.DB
}

// NewRankingRepository создает новый экземпляр rankingRepository
func NewRankingRepository(db *gorm.DB) RankingRepository {
	return &rankingRepository{db: db}
}

func (r *rankingRepository) CreateOrUpdate(profile *domain.RankingProfile) error {
	// Проверяем, существует ли профиль
	var existing domain.RankingProfile
	err := r.db.Where("user_id = ?", profile.UserID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Создаём новый профиль
		return r.db.Create(profile).Error
	}

	// Обновляем существующий
	existing.PairwiseMatrix = profile.PairwiseMatrix
	existing.Weights = profile.Weights
	return r.db.Save(&existing).Error
}

func (r *rankingRepository) GetByUserID(userID uint) (*domain.RankingProfile, error) {
	var profile domain.RankingProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// FavoriteRepository определяет интерфейс для работы с избранным
type FavoriteRepository interface {
	Add(favorite *domain.Favorite) error
	Remove(userID, eventID uint) error
	GetByUserID(userID uint, limit, offset int) ([]domain.Favorite, error)
	IsFavorite(userID, eventID uint) bool
}

// favoriteRepository — реализация FavoriteRepository
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
	return r.db.Where("user_id = ? AND event_id = ?", userID, eventID).
		Delete(&domain.Favorite{}).Error
}

func (r *favoriteRepository) GetByUserID(userID uint, limit, offset int) ([]domain.Favorite, error) {
	var favorites []domain.Favorite
	err := r.db.Where("user_id = ?", userID).
		Preload("Event").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&favorites).Error
	return favorites, err
}

func (r *favoriteRepository) IsFavorite(userID, eventID uint) bool {
	var count int64
	r.db.Model(&domain.Favorite{}).
		Where("user_id = ? AND event_id = ?", userID, eventID).
		Count(&count)
	return count > 0
}
