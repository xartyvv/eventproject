package domain

import (
	"time"

	"gorm.io/gorm"
)

// Event представляет карточку мероприятия с числовыми критериями для МАИ
type Event struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Date        time.Time      `gorm:"not null" json:"date"`
	Location    string         `gorm:"not null" json:"location"`
	Category    string         `json:"category"`
	CreatorID   uint           `gorm:"not null" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Связи
	Creator     *User  `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	FavoritedBy []User `gorm:"many2many:favorites;" json:"-"`

	// 10 критериев для метода анализа иерархий (МАИ)

	// 1. Cost — стоимость билета (в рублях, меньше = лучше)
	Cost float64 `gorm:"default:0" json:"cost"`

	// 2. Duration — длительность мероприятия (в часах)
	Duration float64 `gorm:"default:0" json:"duration"`

	// 3. Capacity — вместимость (кол-во человек)
	Capacity float64 `gorm:"default:0" json:"capacity"`

	// 4. IsWeekend — проходит в выходные (true/false)
	IsWeekend bool `gorm:"default:false" json:"is_weekend"`

	// 5. IsOnline — формат: онлайн (true) или оффлайн (false)
	IsOnline bool `gorm:"default:false" json:"is_online"`

	// 6. AgeRestriction — возрастное ограничение (0 = без ограничений)
	AgeRestriction int `gorm:"default:0" json:"age_restriction"`

	// 7. RequiresRegistration — требуется предварительная регистрация
	RequiresRegistration bool `gorm:"default:false" json:"requires_registration"`

	// 8. OrganizerRating — рейтинг организатора (0-10)
	OrganizerRating float64 `gorm:"default:0" json:"organizer_rating"`

	// 9. TimeOfDay — время проведения: 1=утро, 2=день, 3=вечер, 4=ночь
	TimeOfDay int `gorm:"default:0" json:"time_of_day"`

	// 10. Interactivity — интерактивность (0-10, больше = лучше)
	Interactivity float64 `gorm:"default:0" json:"interactivity"`
}

// TableName указывает имя таблицы для модели Event
func (Event) TableName() string {
	return "events"
}
