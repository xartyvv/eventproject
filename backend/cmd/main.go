package main

import (
	"log"

	"github.com/xartyvv/eventproject/backend/internal/app"
	"github.com/xartyvv/eventproject/backend/internal/config"
	"github.com/xartyvv/eventproject/backend/internal/database"
	"github.com/xartyvv/eventproject/backend/internal/migration"
)

func main() {
	// Загружаем конфигурацию
	config.LoadConfig()

	// Подключаемся к базе данных
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Запускаем миграции (создание таблиц + индексы)
	if err := migration.AutoMigrate(database.DB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Запускаем сервер
	server := app.NewServer()
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
