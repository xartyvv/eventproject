package domain

import (
	"time"

	"gorm.io/gorm"
)

// User представляет зарегистрированного пользователя системы
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"` // "-" не включает в JSON
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Связи
	Events         []Event         `gorm:"foreignKey:CreatorID" json:"events,omitempty"`
	Favorites      []Favorite      `json:"favorites,omitempty"`
	RankingProfile *RankingProfile `gorm:"foreignKey:UserID" json:"ranking_profile,omitempty"`
}

// TableName указывает имя таблицы для модели User
func (User) TableName() string {
	return "users"
}
