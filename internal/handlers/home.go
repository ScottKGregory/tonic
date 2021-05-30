package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scottkgregory/tonic/internal/helpers"
)

type HomeHandler struct {
	Header string
}

func NewHomeHandler(header string) *HomeHandler {
	return &HomeHandler{header}
}

const home = `# Welcome to Tonic!

*To hide this page set {{ .Backtick }}DisableHomepage{{ .Backtick }} to true in Tonic options*

[LOGIN](/auth/login)

---

## Quick start

*You're viewing the site so I assume you've gotten this far already!*

{{ .Backticks }}go

func main() {
	r, authed, err := tonic.Init(cfg.Tonic)
	if err != nil {
		panic(err)
	}

	r.GET("/tonic", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "Hello!"}) })

	authed.GET("/tonic", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Getting here means you're logged in, hey there!"})
	})

	logger := helpers.GetLogger()
	logger.Info().Msg("Starting listener")
	r.Run(":8080")
}

{{ .Backticks }}
`

func (h *HomeHandler) Home() gin.HandlerFunc {
	pageData, err := helpers.MarkdownPage(home, h.Header)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", pageData)
	}
}
