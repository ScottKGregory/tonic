package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/gin-gonic/gin"
	blackfriday "github.com/russross/blackfriday/v2"
)

const htmlHeader = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Tonic | Home</title>
  <link href="//fonts.googleapis.com/css?family=Roboto+Slab:400,300,600" rel="stylesheet" type="text/css">
  <style>
    body {
      font-family: "Roboto Slab", "HelveticaNeue", "Helvetica Neue", Helvetica, Arial, sans-serif;
      text-align: center;
      max-width: 800px;
      margin: 0 auto;
    }
    pre {
      padding: 1em;
      text-align: left;
      border: 1px solid black;
      border-radius: 4px;
    }
  </style>
</head>
<body>
<h1 style="font-size: 5rem">🍸</h1>`

const htmlFooter = "</body></html>"

const home = `# Welcome to Tonic!

*To hide this page set %[1]sDisableHomepage%[1]s to true in Tonic options*

[/login](login)

---

## Quick start

*You're viewing the site so I assume you've gotten this far already!*

%[2]sgo

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

%[2]s
`

func HomePage(c *gin.Context) {
	h := fmt.Sprintf(home, "`", "```")
	bf := blackfriday.WithRenderer(bfchroma.NewRenderer(bfchroma.ChromaOptions(html.TabWidth(2))))
	b := append([]byte(htmlHeader), blackfriday.Run([]byte(h), bf)...)
	c.Data(http.StatusOK, "text/html; charset=utf-8", append(b, []byte(htmlFooter)...))
}

func ErrorPage(override int) gin.HandlerFunc {
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

		content := fmt.Sprintf("# %d - %s", code, status)

		b := append([]byte(htmlHeader), blackfriday.Run([]byte(content), blackfriday.WithRenderer(bfchroma.NewRenderer()))...)
		c.Data(http.StatusOK, "text/html; charset=utf-8", append(b, []byte(htmlFooter)...))
	}
}
