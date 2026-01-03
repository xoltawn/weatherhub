package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xoltawn/weatherhub/docs"
	"github.com/xoltawn/weatherhub/internal/api/handler"
	"github.com/xoltawn/weatherhub/internal/repository"
	weatherrepository "github.com/xoltawn/weatherhub/internal/repository/weather"
	"github.com/xoltawn/weatherhub/internal/service"
	"github.com/xoltawn/weatherhub/pkg/openweathermap"
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

	owmCli := openweathermap.NewOpenWeatherProvider(
		os.Getenv("OPEN_WEATHER_MAP_API_KEY"),
		os.Getenv("OPEN_WEATHER_MAP_BASE_URL"),
		validator.New(),
	)

	cacheTTL, cacheErr := time.ParseDuration(os.Getenv("CACHE_TTL"))
	if cacheErr != nil {
		log.Println("Using default ttl for cache of 1 hour")
		cacheTTL = time.Hour
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,  // use default DB
		PoolSize: 10, // connection pool size
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	weatherRepo := weatherrepository.New(db)
	cachedWeatherRepo := weatherrepository.NewCachedWeatherRepo(weatherRepo, rdb, cacheTTL)
	weatherService := service.NewWeatherService(cachedWeatherRepo, owmCli)

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

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
