package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/xoltawn/weatherhub/internal/domain"
	"github.com/xoltawn/weatherhub/pkg/errutil"
)

type openWeatherProvider struct {
	apiKey    string
	baseURL   string
	validator *validator.Validate
}

func NewOpenWeatherProvider(apiKey, baseURL string, validator *validator.Validate) domain.WeatherProvider {
	return &openWeatherProvider{
		apiKey:    apiKey,
		baseURL:   baseURL,
		validator: validator,
	}
}

type OWMResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description" validate:"required"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Pressure int     `json:"pressure"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Temp     float64 `json:"temp" validate:"required"`
		Humidity int     `json:"humidity" validate:"gte=0,lte=100"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed" validate:"gte=0"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int64 `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country" validate:"required,iso3166_1_alpha2"`
		Sunrise int64   `json:"sunrise"`
		Sunset  int64   `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
	Cod  int    `json:"cod"`
}

func (p *openWeatherProvider) GetForecast(ctx context.Context, city, country string, units domain.Unit) (*domain.WeatherData, error) {
	query := fmt.Sprintf("%s,%s", city, country)
	fullURL := fmt.Sprintf("%s?q=%s&appid=%s&units=%s",
		p.baseURL,
		url.QueryEscape(query),
		p.apiKey,
		units,
	)

	resp, err := http.Get(fullURL)
	if err != nil {
		log.Println(err)

		return nil, errutil.Wrap(domain.ErrThirdParty, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("API error: status %d", resp.StatusCode)

		log.Println(err)

		return nil, errutil.Wrap(domain.ErrThirdParty, err.Error())
	}

	var raw OWMResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		log.Println(err)

		return nil, errutil.Wrap(domain.ErrThirdParty, err.Error())
	}

	if err := p.validator.Struct(raw); err != nil {
		log.Println(err)

		return nil, errutil.Wrap(domain.ErrThirdParty, "weatherapi: provider returned invalid schema")
	}

	return &domain.WeatherData{
		Temperature: raw.Main.Temp,
		Humidity:    raw.Main.Humidity,
		WindSpeed:   raw.Wind.Speed,
		Description: raw.Weather[0].Description,
		CityName:    raw.Name,
		CountryCode: raw.Sys.Country,
	}, nil
}
