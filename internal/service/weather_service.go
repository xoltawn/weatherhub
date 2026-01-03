package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xoltawn/weatherhub/internal/domain"
)

type weatherService struct {
	repo            domain.WeatherRepository
	weatherProvider domain.WeatherProvider
}

func NewWeatherService(repo domain.WeatherRepository, weatherProvider domain.WeatherProvider) domain.WeatherService {
	return &weatherService{
		repo:            repo,
		weatherProvider: weatherProvider,
	}
}

func (s *weatherService) FetchAndStore(ctx context.Context, cityName, country string, units domain.Unit) (*domain.Weather, error) {
	weatherApiResp, err := s.weatherProvider.GetForecast(ctx, strings.ToLower(cityName), strings.ToLower(country), units)
	if err != nil {
		return nil, err
	}

	weather := &domain.Weather{
		ID:          uuid.New(),
		CityName:    cityName,
		Country:     country,
		Temperature: weatherApiResp.Temperature,
		Description: weatherApiResp.Description,
		Humidity:    weatherApiResp.Humidity,
		WindSpeed:   weatherApiResp.WindSpeed,
		FetchedAt:   time.Now(),
		Unit:        units,
	}

	if err := s.repo.Create(ctx, weather); err != nil {
		return nil, err
	}

	return weather, nil
}

func (s *weatherService) GetAllRecords(ctx context.Context) ([]domain.Weather, error) {
	return s.repo.GetAll(ctx)
}

func (s *weatherService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Weather, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *weatherService) UpdateRecord(ctx context.Context, id uuid.UUID, updates *domain.Weather) (*domain.Weather, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Temperature = updates.Temperature
	existing.Description = updates.Description
	existing.Humidity = updates.Humidity
	existing.WindSpeed = updates.WindSpeed
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *weatherService) DeleteRecord(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *weatherService) GetLatest(ctx context.Context, cityName string) (*domain.Weather, error) {
	return s.repo.GetLatestByCity(ctx, strings.ToLower(cityName))
}
