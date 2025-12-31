package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/xoltawn/weatherhub/internal/domain"
)

type weatherService struct {
	repo domain.WeatherRepository
}

func NewWeatherService(repo domain.WeatherRepository) domain.WeatherService {
	return &weatherService{
		repo: repo,
	}
}

func (s *weatherService) FetchAndStore(ctx context.Context, cityName, country string) (*domain.Weather, error) {
	// TODO: fetch data from openweathermap

	weather := &domain.Weather{
		ID:          uuid.New(),
		CityName:    cityName,
		Country:     country,
		Temperature: 0,  //TODO
		Description: "", //TODO
		Humidity:    0,  //TODO
		WindSpeed:   0,  //TODO
		FetchedAt:   time.Now(),
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
	return s.repo.GetLatestByCity(ctx, cityName)
}
