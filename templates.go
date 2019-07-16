package main

import "text/template"

var indexTmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
  <head>
  <title>{{.Title}}</title>
  <link rel='icon' href='/favicon.ico' type='image/x-icon'/ >
    <style>
    ul {
      display: flex;
      flex-wrap: wrap;
      padding: 0;
      list-style-type: none;
    }
    li { flex: 0 0 33%; }
    li { text-align: center; }
    li:nth-child(n) { background-color: #eed; }
    li:nth-child(6n+4) { background-color: lightgray; }
    li:nth-child(6n+5) { background-color: lightgray; }
    li:nth-child(6n+6) { background-color: lightgray; }
    </style>
  </head>
  <body>
    <h1><a href="https://{{.Host}}/">{{.Title}}</a></h1>
    {{range .Links -}}
    <h3><a href="{{.URL}}">{{.Title}}</a></h3>
    {{end -}}
    <ul>
    {{range .Paths}}{{if ne .Repo ""}}  <li>{{.Path}}</li><li><a href="https://godoc.org/{{$.Host}}{{.Path}}">GoDoc</a></li><li><a href="{{.Repo}}">Code</a></li>
    {{end}}{{end}}</ul>
{{if ne .Src "" -}}
    (<a href="{{.Src}}">source</a>)
{{end -}}
  </body>
</html>
`))

var vanityTmpl = template.Must(template.New("vanity").Parse(`<!DOCTYPE html>
<html>
  <head>
    <link rel='icon' href='/favicon.ico' type='image/x-icon'/ >
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <meta name="go-import" content="{{.Host}}{{.Path}} {{.VCS}} {{.Repo}}">
    <meta name="go-source" content="{{.Host}}{{.Path}} {{.Display}}">
    <meta http-equiv="refresh" content="0; url=https://godoc.org/{{.Host}}{{.Path}}/{{.Subpath}}">
  </head>
  <body>
    Nothing to see here; <a href="https://godoc.org/{{.Host}}{{.Path}}/{{.Subpath}}">See the package on godoc</a>.
  </body>
</html>
`))
