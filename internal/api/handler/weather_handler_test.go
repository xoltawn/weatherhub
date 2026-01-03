package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xoltawn/weatherhub/internal/api/handler"
	"github.com/xoltawn/weatherhub/internal/domain"
	"github.com/xoltawn/weatherhub/internal/repository/mocks"
	"github.com/xoltawn/weatherhub/internal/service"
)

func TestWeatherHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := mocks.NewWeatherRepository(t)

	weatherService := service.NewWeatherService(mockRepo, nil)

	h := handler.NewWeatherHandler(weatherService)

	router := gin.Default()
	router.GET("/weather/:id", h.GetByID)

	t.Run("success", func(t *testing.T) {
		// Arrange
		id := uuid.New()
		expectedWeather := &domain.Weather{
			ID:          id,
			CityName:    "London",
			Temperature: 20,
		}

		mockRepo.On("GetByID", mock.Anything, id).Return(expectedWeather, nil)

		// Act
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/weather/"+id.String(), nil)
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var resp domain.Weather
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "London", resp.CityName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid-uuid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/weather/not-a-uuid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
