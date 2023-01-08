package main

import (
	_ "embed"
	"html/template"
)

//go:embed index.html
var mainTemplateStr string

var t = template.Must(template.New("main").Parse(mainTemplateStr))
