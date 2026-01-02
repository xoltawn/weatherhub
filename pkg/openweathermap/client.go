package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/xoltawn/weatherhub/internal/domain"
	"github.com/xoltawn/weatherhub/pkg/errutil"
)

type openWeatherProvider struct {
	apiKey  string
	baseURL string
}

func NewOpenWeatherProvider(apiKey, baseURL string) domain.WeatherProvider {
	return &openWeatherProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
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
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
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
		Country string  `json:"country"`
		Sunrise int64   `json:"sunrise"`
		Sunset  int64   `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
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
		return nil, errutil.Wrap(domain.ErrThirdParty, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errutil.Wrap(domain.ErrThirdParty, fmt.Errorf("API error: status %d", resp.StatusCode).Error())
	}

	var raw OWMResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, errutil.Wrap(domain.ErrThirdParty, err.Error())
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
