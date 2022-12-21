package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/pkg/api/errors"
	"github.com/scottkgregory/tonic/pkg/dependencies"
)

type ResponseModel struct {
	Data       interface{}       `json:"data,omitempty"`
	Error      string            `json:"error,omitempty"`
	Validation map[string]string `json:"validation,omitempty"`
}

func SmartResponse(c *gin.Context, data interface{}, err error) {
	if errors.Is(err, &errors.ValidationErr{}) {
		ValidationErrorResponse(c, err.(*errors.ValidationErr))
		return
	} else if errors.Is(err, &errors.NotFoundErr{}) {
		NotFoundResponse(c, err.(*errors.NotFoundErr))
		return
	} else if errors.Is(err, &errors.UnauthorisedErr{}) {
		UnauthorisedResponse(c, err.(*errors.UnauthorisedErr))
		return
	} else if errors.Is(err, &errors.ForbiddenErr{}) {
		ForbiddenResponse(c, err.(*errors.ForbiddenErr))
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

func UnauthorisedResponse(c *gin.Context, errs ...*errors.UnauthorisedErr) {
	err := errors.NewUnauthorisedError()
	if len(errs) == 1 {
		err = errs[0]
	}

	ErrorResponse(c, http.StatusUnauthorized, err)
}

func ForbiddenResponse(c *gin.Context, errs ...*errors.ForbiddenErr) {
	err := errors.NewForbiddenError()
	if len(errs) == 1 {
		err = errs[0]
	}

	ErrorResponse(c, http.StatusForbidden, err)
}

func NotFoundResponse(c *gin.Context, errs ...*errors.NotFoundErr) {
	err := errors.NewNotFoundError("")
	if len(errs) == 1 {
		err = errs[0]
	}

	ErrorResponse(c, http.StatusNotFound, err)
}

func ValidationErrorResponse(c *gin.Context, errs ...*errors.ValidationErr) {
	err := errors.NewValidationError()
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

	if e, ok := err.(errors.TonicError); ok {
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
