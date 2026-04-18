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
	Creator   *User      `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	FavoritedBy []User   `gorm:"many2many:favorites;" json:"-"`

	// 15 критериев для метода анализа иерархий (МАИ)
	// Все критерии нормализованы в диапазон 0-1 или имеют числовое значение

	// 1. Cost — стоимость билета (в рублях, меньше = лучше)
	Cost float64 `gorm:"default:0" json:"cost"`

	// 2. Distance — удалённость от центра (в км, меньше = лучше)
	Distance float64 `gorm:"default:0" json:"distance"`

	// 3. Duration — длительность мероприятия (в часах)
	Duration float64 `gorm:"default:0" json:"duration"`

	// 4. Rating — рейтинг мероприятия (0-10, больше = лучше)
	Rating float64 `gorm:"default:0" json:"rating"`

	// 5. Capacity — вместимость (кол-во человек)
	Capacity float64 `gorm:"default:0" json:"capacity"`

	// 6. IsWeekend — проходит в выходные (1 = да, 0 = нет)
	IsWeekend bool `gorm:"default:false" json:"is_weekend"`

	// 7. IsOnline — формат: онлайн (true) или оффлайн (false)
	IsOnline bool `gorm:"default:false" json:"is_online"`

	// 8. AgeRestriction — возрастное ограничение (0 = без ограничений)
	AgeRestriction int `gorm:"default:0" json:"age_restriction"`

	// 9. RequiresRegistration — требуется предварительная регистрация
	RequiresRegistration bool `gorm:"default:false" json:"requires_registration"`

	// 10. OrganizerRating — рейтинг организатора (0-10)
	OrganizerRating float64 `gorm:"default:0" json:"organizer_rating"`

	// 11. IsFree — бесплатное мероприятие (true/false)
	IsFree bool `gorm:"default:false" json:"is_free"`

	// 12. TimeOfDay — время проведения: 1=утро, 2=день, 3=вечер, 4=ночь
	TimeOfDay int `gorm:"default:0" json:"time_of_day"`

	// 13. Accessibility — доступность для маломобильных (0-10)
	Accessibility float64 `gorm:"default:0" json:"accessibility"`

	// 14. Popularity — популярность мероприятия (кол-во участников)
	Popularity float64 `gorm:"default:0" json:"popularity"`

	// 15. Interactivity — интерактивность (0-10, больше = лучше)
	Interactivity float64 `gorm:"default:0" json:"interactivity"`
}

// TableName указывает имя таблицы для модели Event
func (Event) TableName() string {
	return "events"
}
