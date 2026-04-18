package migration

import (
	"fmt"
	"log"

	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"gorm.io/gorm"
)

// AutoMigrate запускает миграцию всех таблиц с полной настройкой индексов и foreign keys
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Включаем foreign keys для SQLite (если потребуется) и для PostgreSQL
	db.Exec("SET session_replication_role = 'origin';")

	// 1. Создаём таблицы
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Event{},
		&domain.RankingProfile{},
		&domain.Favorite{},
	); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	// 2. Добавляем дополнительные индексы вручную
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("create indexes failed: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// createIndexes создаёт дополнительные индексы для оптимизации запросов
func createIndexes(db *gorm.DB) error {
	indexes := []struct {
		table  string
		name   string
		column string
	}{
		{"events", "idx_events_category", "category"},
		{"events", "idx_events_date", "date"},
		{"events", "idx_events_creator", "creator_id"},
		{"events", "idx_events_category_date", "category, date"},
		{"favorites", "idx_favorites_user", "user_id"},
		{"ranking_profiles", "idx_ranking_profiles_user", "user_id"},
	}

	for _, idx := range indexes {
		// Проверяем, существует ли индекс
		var count int64
		db.Raw(`
			SELECT COUNT(*) FROM pg_indexes 
			WHERE tablename = ? AND indexname = ?
		`, idx.table, idx.name).Scan(&count)

		if count == 0 {
			sql := fmt.Sprintf("CREATE INDEX %s ON %s (%s)", idx.name, idx.table, idx.column)
			if err := db.Exec(sql).Error; err != nil {
				log.Printf("Warning: failed to create index %s: %v", idx.name, err)
			} else {
				log.Printf("Created index: %s on %s(%s)", idx.name, idx.table, idx.column)
			}
		} else {
			log.Printf("Index already exists: %s", idx.name)
		}
	}

	return nil
}
