package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Unit string

const (
	Metric   Unit = "metric"
	Imperial Unit = "imperial"
)

type Weather struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CityName    string    `json:"city_name"`
	Country     string    `json:"country"`
	Temperature float64   `json:"temperature"`
	Unit        Unit      `json:"unit"`
	Description string    `json:"description"`
	Humidity    int       `json:"humidity"`
	WindSpeed   float64   `json:"wind_speed"`
	FetchedAt   time.Time `json:"fetched_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//go:generate mockery --name=WeatherRepository --output=../repository/mocks --case=underscore
type WeatherRepository interface {
	Create(ctx context.Context, weather *Weather) error
	GetAll(ctx context.Context) ([]Weather, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Weather, error)
	GetLatestByCity(ctx context.Context, cityName string) (*Weather, error)
	Update(ctx context.Context, weather *Weather) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type WeatherService interface {
	FetchAndStore(ctx context.Context, cityName, country string, units Unit) (*Weather, error)
	GetAllRecords(ctx context.Context) ([]Weather, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Weather, error)
	UpdateRecord(ctx context.Context, id uuid.UUID, updates *Weather) (*Weather, error)
	DeleteRecord(ctx context.Context, id uuid.UUID) error
	GetLatest(ctx context.Context, cityName string) (*Weather, error)
}

type WeatherData struct {
	Temperature float64
	Humidity    int
	WindSpeed   float64
	Description string
	CityName    string
	CountryCode string
}

type WeatherProvider interface {
	GetForecast(ctx context.Context, city, country string, units Unit) (*WeatherData, error)
}
