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
