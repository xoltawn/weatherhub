package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xoltawn/weatherhub/internal/domain"
)

func RespondWithError(c *gin.Context, err error) {
	println("LOGGING ERROR:", err.Error())

	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "The requested weather data was not found."})
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please check your city name or country code."})
	case errors.Is(err, domain.ErrThirdParty):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Weather provider is currently down. Try again later."})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred on our end."})
	}
}
