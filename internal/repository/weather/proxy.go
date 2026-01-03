package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/xoltawn/weatherhub/internal/domain"
)

type cachedWeatherRepo struct {
	realRepo domain.WeatherRepository
	redis    *redis.Client
	ttl      time.Duration
}

func NewCachedWeatherRepo(real domain.WeatherRepository, rdb *redis.Client, ttl time.Duration) domain.WeatherRepository {
	return &cachedWeatherRepo{
		realRepo: real,
		redis:    rdb,
		ttl:      ttl,
	}
}

func (r *cachedWeatherRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Weather, error) {
	cacheKey := r.fmtKey(id)

	val, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var weather domain.Weather
		if json.Unmarshal([]byte(val), &weather) == nil {
			return &weather, nil
		}
	}

	weather, err := r.realRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data, marshalErr := json.Marshal(weather)
	if marshalErr != nil {
		//log it or put in job , there is no need to affect user response for this (not sensitive data)
		return weather, nil
	}

	setCmd := r.redis.Set(ctx, cacheKey, data, r.ttl)
	if setCmd.Err() != nil {
		//JUST log , there is no need to affect user response for this
	}

	return weather, nil
}

func (r *cachedWeatherRepo) Create(ctx context.Context, w *domain.Weather) error {
	if err := r.realRepo.Create(ctx, w); err != nil {
		return err
	}

	data, marshalErr := json.Marshal(w)
	if marshalErr != nil {
		//log it or put in job , there is no need to affect user response for this (not sensitive data)
		return nil
	}

	setCmd := r.redis.Set(ctx, "weather:"+w.ID.String(), data, r.ttl)
	if setCmd.Err() != nil {
		//JUST log , there is no need to affect user response for this
	}

	return nil
}

func (r *cachedWeatherRepo) fmtKey(id uuid.UUID) string {
	return fmt.Sprintf("weather:%s", id.String())
}

func (r *cachedWeatherRepo) Update(ctx context.Context, w *domain.Weather) error {
	if err := r.realRepo.Update(ctx, w); err != nil {
		return err
	}

	data, marshalErr := json.Marshal(w)
	if marshalErr != nil {
		//log it or put in job , there is no need to affect user response for this (not sensitive data)
		return nil
	}

	setCmd := r.redis.Set(ctx, r.fmtKey(w.ID), data, r.ttl)
	if setCmd.Err() != nil {
		//log it or put in job , there is no need to affect user response for this (not sensitive data)
	}

	return nil
}

func (r *cachedWeatherRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.realRepo.Delete(ctx, id); err != nil {
		return err
	}

	delCmd := r.redis.Del(ctx, r.fmtKey(id))
	if delCmd.Err() != nil {
		//log it or put in job , there is no need to affect user response for this (not sensitive data)
	}

	return nil
}

func (r *cachedWeatherRepo) GetAll(ctx context.Context) ([]domain.Weather, error) {
	return r.realRepo.GetAll(ctx)
}

func (r *cachedWeatherRepo) GetLatestByCity(ctx context.Context, cityName string) (*domain.Weather, error) {
	return r.realRepo.GetLatestByCity(ctx, cityName)
}
