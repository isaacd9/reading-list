package main

import (
	"text/template"
)

const tmpl = `{{- range $date, $articles := . -}}
- {{ $date }}
	{{- range $article := $articles}}
	- [{{ $article.Title }}]({{ $article.URL }})
	{{- end }}
{{ end }}
`

var markdown = template.Must(template.New("reading-list").Parse(tmpl))
