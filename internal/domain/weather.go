package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Weather struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CityName    string    `json:"city_name"`
	Country     string    `json:"country"`
	Temperature float64   `json:"temperature"`
	Description string    `json:"description"`
	Humidity    int       `json:"humidity"`
	WindSpeed   float64   `json:"wind_speed"`
	FetchedAt   time.Time `json:"fetched_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WeatherRepository interface {
	Create(ctx context.Context, weather *Weather) error
	GetAll(ctx context.Context) ([]Weather, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Weather, error)
	GetLatestByCity(ctx context.Context, cityName string) (*Weather, error)
	Update(ctx context.Context, weather *Weather) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type WeatherService interface {
	FetchAndStore(ctx context.Context, cityName, country string) (*Weather, error)
	GetAllRecords(ctx context.Context) ([]Weather, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Weather, error)
	UpdateRecord(ctx context.Context, id uuid.UUID, updates *Weather) (*Weather, error)
	DeleteRecord(ctx context.Context, id uuid.UUID) error
	GetLatest(ctx context.Context, cityName string) (*Weather, error)
}
