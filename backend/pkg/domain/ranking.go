package domain

import (
	"time"

	"gorm.io/gorm"
)

// CriterionName — список названий критериев для МАИ
var CriterionNames = []string{
	"Стоимость билета",
	"Длительность",
	"Вместимость",
	"Выходной день",
	"Онлайн формат",
	"Возрастное ограничение",
	"Требуется регистрация",
	"Рейтинг организатора",
	"Время проведения",
	"Интерактивность",
}

// Matrix10x10 представляет матрицу парных сравнений 10x10 (для 10 критериев)
type Matrix10x10 [10][10]float64

// RankingProfile хранит профиль ранжирования пользователя (матрицу парных сравнений)
type RankingProfile struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Матрица парных сравнений 15x15 (JSON encoded)
	// PairwiseMatrix[i][j] = насколько критерий i важнее критерия j (шкала 1-9)
	PairwiseMatrix string `gorm:"type:text" json:"pairwise_matrix"`

	// Вычисленные веса критериев (JSON encoded)
	// Weights[i] — вес i-го критерия (нормализованный, сумма = 1)
	Weights string `gorm:"type:text" json:"weights"`

	// Связи
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName указывает имя таблицы для модели RankingProfile
func (RankingProfile) TableName() string {
	return "ranking_profiles"
}

// RankedEvent — результат ранжирования мероприятия для конкретного пользователя
type RankedEvent struct {
	EventID     uint      `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	Location    string    `json:"location"`
	Category    string    `json:"category"`
	Score       float64   `json:"score"`    // Итоговый балл по МАИ
	Rank        int       `json:"rank"`     // Позиция в рейтинге
	Criteria    []float64 `json:"criteria"` // Значения всех 15 критериев
}
