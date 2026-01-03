package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/xoltawn/weatherhub/internal/domain"
)

func RespondWithError(c *gin.Context, err error) {
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		translated := vErrs.Translate(trans)
		var details []map[string]string

		for field, msg := range translated {
			details = append(details, map[string]string{
				"field":   field,
				"message": msg,
			})
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": details,
		})
		return
	}

	var statusCode int
	var message string

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
		message = "The requested resource was not found."
	case errors.Is(err, domain.ErrThirdParty):
		statusCode = http.StatusServiceUnavailable
		message = "External weather service is unavailable."
	default:
		statusCode = http.StatusInternalServerError
		message = "An internal server error occurred."
	}

	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

var (
	uni   *ut.UniversalTranslator
	trans ut.Translator
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		enLocale := en.New()
		uni = ut.New(enLocale, enLocale)
		trans, _ = uni.GetTranslator("en")

		en_translations.RegisterDefaultTranslations(v, trans)
	}
}

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func MapValidationErrors(err error) (bool, []ValidationErrorResponse) {
	var errs []ValidationErrorResponse

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return false, errs
	}

	translated := ve.Translate(trans)

	for field, message := range translated {
		errs = append(errs, ValidationErrorResponse{
			Field:   field,
			Message: message,
		})
	}

	return true, errs
}
