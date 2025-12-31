package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xoltawn/weatherhub/internal/domain"
)

type WeatherHandler struct {
	service domain.WeatherService
}

func NewWeatherHandler(service domain.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: service}
}

func (h *WeatherHandler) RegisterRoutes(rg *gin.RouterGroup) {
	weather := rg.Group("/weather")
	{
		weather.GET("", h.GetAll)
		weather.GET("/:id", h.GetByID)
		weather.POST("", h.Create)
		weather.PUT("/:id", h.Update)
		weather.DELETE("/:id", h.Delete)
		weather.GET("/latest/:cityName", h.GetLatest)
	}
}

// Create godoc
// @Summary      Fetch and store weather
// @Description  Calls OpenWeatherMap API for a city/country and saves the result to the database
// @Tags         weather
// @Accept       json
// @Produce      json
// @Param        request  body      object{cityName=string,country=string}  true  "City and Country codes"
// @Success      201      {object}  domain.Weather
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /weather [post]
func (h *WeatherHandler) Create(c *gin.Context) {
	var input struct {
		CityName string `json:"cityName" binding:"required"`
		Country  string `json:"country" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result, err := h.service.FetchAndStore(c.Request.Context(), input.CityName, input.Country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetAll godoc
// @Summary      List all weather records
// @Description  Retrieve every weather record currently stored in the database
// @Tags         weather
// @Produce      json
// @Success      200  {array}   domain.Weather
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /weather [get]
func (h *WeatherHandler) GetAll(c *gin.Context) {
	weathers, err := h.service.GetAllRecords(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, weathers)
}

// GetByID godoc
// @Summary      Get weather by ID
// @Description  Retrieve a specific weather record using its UUID
// @Tags         weather
// @Produce      json
// @Param        id   path      string  true  "Weather UUID"
// @Success      200  {object}  domain.Weather
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /weather/{id} [get]
func (h *WeatherHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	weather, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, weather)
}

// Update godoc
// @Summary      Update a weather record
// @Description  Modify fields of an existing weather record by ID
// @Tags         weather
// @Accept       json
// @Produce      json
// @Param        id       path      string          true  "Weather UUID"
// @Param        updates  body      domain.Weather  true  "Fields to update"
// @Success      200      {object}  domain.Weather
// @Failure      400      {object}  map[string]string
// @Security     BearerAuth
// @Router       /weather/{id} [put]
func (h *WeatherHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var updates domain.Weather
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.UpdateRecord(c.Request.Context(), id, &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Delete godoc
// @Summary      Delete a weather record
// @Description  Remove a weather record from the database by ID
// @Tags         weather
// @Param        id   path      string  true  "Weather UUID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Security     BearerAuth
// @Router       /weather/{id} [delete]
func (h *WeatherHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := h.service.DeleteRecord(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}

// GetLatest godoc
// @Summary      Get latest city weather
// @Description  Retrieve the most recently fetched weather record for a specific city
// @Tags         weather
// @Produce      json
// @Param        cityName  path      string  true  "City Name"
// @Success      200       {object}  domain.Weather
// @Failure      404       {object}  map[string]string
// @Security     BearerAuth
// @Router       /weather/latest/{cityName} [get]
func (h *WeatherHandler) GetLatest(c *gin.Context) {
	cityName := c.Param("cityName")
	result, err := h.service.GetLatest(c.Request.Context(), cityName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No records found for this city"})
		return
	}
	c.JSON(http.StatusOK, result)
}
