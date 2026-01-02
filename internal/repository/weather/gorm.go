package weather

import (
	"context"

	"github.com/google/uuid"
	"github.com/xoltawn/weatherhub/internal/domain"
	"github.com/xoltawn/weatherhub/internal/repository"
	"gorm.io/gorm"
)

type weatherRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.WeatherRepository {
	return &weatherRepo{db: db}
}

func (r *weatherRepo) Create(ctx context.Context, weather *domain.Weather) error {
	err := r.db.
		WithContext(ctx).
		Create(weather).Error
	if err != nil {
		return repository.MapGormError(err, "repository.Weather.Create")
	}

	return err
}

func (r *weatherRepo) GetAll(ctx context.Context) ([]domain.Weather, error) {
	var records []domain.Weather

	err := r.db.
		WithContext(ctx).
		Find(&records).Error
	if err != nil {
		return nil, repository.MapGormError(err, "repository.Weather.GetAll")
	}

	return records, err
}

func (r *weatherRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Weather, error) {
	var weather domain.Weather

	err := r.db.
		WithContext(ctx).
		Take(&weather, "id = ?", id).
		Error
	if err != nil {
		return nil, repository.MapGormError(err, "repository.Weather.GetByID")
	}

	return &weather, nil
}

func (r *weatherRepo) GetLatestByCity(ctx context.Context, cityName string) (*domain.Weather, error) {
	var weather domain.Weather

	err := r.db.
		WithContext(ctx).
		Where("city_name = ?", cityName).
		Order("fetched_at DESC").
		First(&weather).Error
	if err != nil {
		return nil, repository.MapGormError(err, "repository.Weather.GetLatestByCity")
	}

	return &weather, nil
}

func (r *weatherRepo) Update(ctx context.Context, weather *domain.Weather) error {
	err := r.db.
		WithContext(ctx).
		Save(weather).Error
	if err != nil {
		return repository.MapGormError(err, "repository.Weather.Update")
	}

	return nil
}

func (r *weatherRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.
		WithContext(ctx).
		Delete(&domain.Weather{}, "id = ?", id).
		Error
	if err != nil {
		return repository.MapGormError(err, "repository.Weather.Delete")
	}

	return nil
}
