package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xoltawn/weatherhub/docs"
	"github.com/xoltawn/weatherhub/internal/api/handler"
	"github.com/xoltawn/weatherhub/internal/repository"
	weatherrepository "github.com/xoltawn/weatherhub/internal/repository/weather"
	"github.com/xoltawn/weatherhub/internal/service"
)

// @title WeatherHub API
// @version 1.0
// @description This is a weather data server.
// @BasePath /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := repository.InitDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	weatherRepo := weatherrepository.New(db)
	weatherService := service.NewWeatherService(weatherRepo)

	router := gin.Default()
	api := router.Group("/api/v1")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	weatherHandler := handler.NewWeatherHandler(weatherService)
	weatherHandler.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
