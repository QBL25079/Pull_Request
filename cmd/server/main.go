package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"pr-reviewer-service/internal/config"
	"pr-reviewer-service/internal/handler"
	"pr-reviewer-service/internal/repository"
	"pr-reviewer-service/internal/service"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg, _ := config.Load()

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL!")

	repo := repository.NewPostgresRepository(db)
	userService := service.New(repo)
	h := handler.NewHandler(userService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
	v1 := e.Group("/api/v1")

	v1.POST("/teams", h.CreateTeam)
	v1.PUT("/teams/:teamName", h.UpdateTeamName)

	v1.POST("/users", h.CreateUser)
	v1.GET("/users/:userID", h.GetUserByID)

	v1.POST("/pull-requests", h.CreatePullRequest)
	v1.GET("/pull-requests/:prID", h.GetPullRequestByID)

	log.Printf("Server starting on :%s", cfg.HTTPPort)
	log.Fatal(e.Start(":" + cfg.HTTPPort))
}
