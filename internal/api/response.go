package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	tonicErrors "github.com/scottkgregory/tonic/internal/api/errors"
	"github.com/scottkgregory/tonic/internal/dependencies"
)

type ResponseModel struct {
	Data       interface{}       `json:"data,omitempty"`
	Error      string            `json:"error,omitempty"`
	Validation map[string]string `json:"validation,omitempty"`
}

func SmartResponse(c *gin.Context, data interface{}, err error) {
	if errors.Is(err, &tonicErrors.ValidationErr{}) {
		ValidationErrorResponse(c, err.(*tonicErrors.ValidationErr))
		return
	} else if errors.Is(err, &tonicErrors.NotFoundErr{}) {
		NotFoundResponse(c, err.(*tonicErrors.NotFoundErr))
		return
	} else if errors.Is(err, &tonicErrors.UnauthorisedErr{}) {
		UnauthorisedResponse(c, err.(*tonicErrors.UnauthorisedErr))
		return
	} else if errors.Is(err, &tonicErrors.ForbiddenErr{}) {
		ForbiddenResponse(c, err.(*tonicErrors.ForbiddenErr))
		return
	} else if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if data == nil {
		NoContentResponse(c)
		return
	}

	SuccessResponse(c, data)
}

func NoContentResponse(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNoContent)
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseModel{Data: data})
}

func UnauthorisedResponse(c *gin.Context, errs ...*tonicErrors.UnauthorisedErr) {
	err := tonicErrors.NewUnauthorisedError()
	if len(errs) == 1 {
		err = errs[0]
	}

	ErrorResponse(c, http.StatusUnauthorized, err)
}

func ForbiddenResponse(c *gin.Context, errs ...*tonicErrors.ForbiddenErr) {
	err := tonicErrors.NewForbiddenError()
	if len(errs) == 1 {
		err = errs[0]
	}

	ErrorResponse(c, http.StatusForbidden, err)
}

func NotFoundResponse(c *gin.Context, errs ...*tonicErrors.NotFoundErr) {
	err := tonicErrors.NewNotFoundError("")
	if len(errs) == 1 {
		err = errs[0]
	}

	ErrorResponse(c, http.StatusNotFound, err)
}

func ValidationErrorResponse(c *gin.Context, errs ...*tonicErrors.ValidationErr) {
	err := tonicErrors.NewValidationError()
	if len(errs) == 1 {
		err = errs[0]
	}

	ErrorResponse(c, http.StatusBadRequest, err, err.Validation)
}

func ErrorResponse(c *gin.Context, code int, err error, validation ...map[string]string) {
	dependencies.GetLogger(c).Err(err).Msg("Error processing request")
	var val map[string]string
	if len(validation) == 1 {
		val = validation[0]
	}

	if e, ok := err.(tonicErrors.TonicError); ok {
		c.JSON(code, &ResponseModel{
			Error:      e.External(),
			Validation: val,
		})

		return
	}

	c.JSON(code, &ResponseModel{
		Error:      err.Error(),
		Validation: val,
	})
}
