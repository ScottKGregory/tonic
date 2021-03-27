package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseModel struct {
	Data       interface{}       `json:"data,omitempty"`
	Error      string            `json:"error,omitempty"`
	Validation map[string]string `json:"validation,omitempty"`
}

func APISuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseModel{Data: data})
}

func APIErrorResponse(c *gin.Context, code int, err string, validation ...map[string]string) {
	var val map[string]string
	if len(validation) == 1 {
		val = validation[0]
	}

	c.JSON(code, &ResponseModel{
		Error:      err,
		Validation: val,
	})
}
