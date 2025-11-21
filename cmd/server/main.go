package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"pr-reviewer-service/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	DB, err := initDb(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	defer DB.Close()

	if err := runMigrations(cfg.GetDSN()); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("health", func(c echo.Context) error {
		if err := DB.Ping(); err != nil {
			return c.String(http.StatusInternalServerError, "Ошибка подключения базы данных")
		}
		return c.String(http.StatusOK, "Сервер успешно запущен")
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to PR Reviewer Service API.")
	})

	log.Printf("Сервер запускается на порту %s...", cfg.HTTPPort)
	e.Logger.Fatal(e.Start(":" + cfg.HTTPPort))
}

func initDb(cfg *config.Config) (*sql.DB, error) {
	const maxRetries = 10

	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", cfg.GetDSN())
		if err == nil {
			err = db.Ping()
			if err == nil {
				db.SetMaxOpenConns(25)
				db.SetMaxIdleConns(10)
				log.Println("Успешное подключение к PostgreSQL.")
				return db, nil
			}
		}
		log.Printf("Не удалось подключиться к бд (попытка %d,%d)", i+1, maxRetries)
		time.Sleep(time.Second * 3)
	}
	return nil, fmt.Errorf("не удалось подключиться к БД после %d попыток: %w", maxRetries, err)
}

func runMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("Не удалось создать экземпляр мигратора, %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка при выполнении миграций Up: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("Миграции не требуются, схема актуальна.")
	}
	return nil
}
