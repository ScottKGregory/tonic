package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/internal/constants"
	"github.com/scottkgregory/tonic/internal/helpers"
)

const errMarkdown = `
# %d - %s

**Request ID:** %s

[HOME](/)
`

type ErrorHandler struct {
	Header string
}

func NewErrorHandler(header string) *ErrorHandler {
	return &ErrorHandler{header}
}

func (h *ErrorHandler) Error(override int) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 0
		if override != 0 {
			code = override
		} else {
			var err error
			code, err = strconv.Atoi(c.Param("code"))
			if err != nil {
				code = http.StatusInternalServerError
			}
		}

		status := http.StatusText(code)
		if status == "" {
			code = http.StatusInternalServerError
			status = http.StatusText(http.StatusInternalServerError)
		}

		msg := fmt.Sprintf(errMarkdown, code, status, c.GetString(constants.RequestIDKey))
		pageData, err := helpers.MarkdownPage(msg, h.Header)
		if err != nil {
			c.String(code, msg)
		}

		c.Data(code, "text/html; charset=utf-8", pageData)
	}
}
