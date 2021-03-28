package helpers

import (
	"bytes"
	"html/template"

	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
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
		a {
			font-weight: 600;
			font-size: 1.25em;
			text-decoration: none;
		}
  </style>
</head>
<body>
<h1 style="font-size: 5rem">🍸</h1>`

const htmlFooter = "</body></html>"

type vars struct {
	Backticks string
	Backtick  string
}

func MarkdownPage(md string) ([]byte, error) {
	tmpl, err := template.New("md").Parse(md)
	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, vars{
		Backticks: "```",
		Backtick:  "`",
	})
	if err != nil {
		return nil, err
	}

	bf := blackfriday.WithRenderer(bfchroma.NewRenderer(bfchroma.ChromaOptions(html.TabWidth(2))))
	b := append([]byte(htmlHeader), blackfriday.Run(tpl.Bytes(), bf)...)
	return append(b, []byte(htmlFooter)...), nil
}
