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
<h1 style="font-size: 5rem">{{ .Header }}</h1>`

const htmlFooter = "</body></html>"

type vars struct {
	Backticks string
	Backtick  string
	Header    string
}

func MarkdownPage(md, pageHeader string) ([]byte, error) {
	v := vars{
		Backticks: "```",
		Backtick:  "`",
		Header:    pageHeader,
	}

	headerTemplate, err := template.New("header").Parse(htmlHeader)
	if err != nil {
		return nil, err
	}

	var headerBytes bytes.Buffer
	err = headerTemplate.Execute(&headerBytes, v)
	if err != nil {
		return nil, err
	}

	markdownTemplate, err := template.New("md").Parse(md)
	if err != nil {
		return nil, err
	}

	var markdownBytes bytes.Buffer
	err = markdownTemplate.Execute(&markdownBytes, v)
	if err != nil {
		return nil, err
	}

	bf := blackfriday.WithRenderer(bfchroma.NewRenderer(bfchroma.ChromaOptions(html.TabWidth(2))))
	b := append(headerBytes.Bytes(), blackfriday.Run(markdownBytes.Bytes(), bf)...)
	return append(b, []byte(htmlFooter)...), nil
}
