package domain

import (
	"time"

	"gorm.io/gorm"
)

// Favorite представляет связь пользователя с избранным мероприятием
type Favorite struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	EventID   uint           `gorm:"not null" json:"event_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Связи
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Event *Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
}

// TableName указывает имя таблицы для модели Favorite
func (Favorite) TableName() string {
	return "favorites"
}
