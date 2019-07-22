package main

import "text/template"

// This is the index page. TODO: prettify it.
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
{{- range .Paths}}  {{if and (ne .Repo "") (ne .Wildcard true)}}
      <li><a href="{{.Path}}">{{.Path}}</a></li>
      <li><a href="https://godoc.org/{{$.Host}}{{.Path}}">GoDoc</a></li>
      <li><a href="{{.Repo}}">Code</a></li>
{{end}}{{- end -}}</ul>
{{if ne .Src ""}}    (<a href="{{.Src}}">source</a>){{end}}
  </body>
</html>
`))

// This is used for requests where go-get=1 is present.
var gogetTmpl = template.Must(template.New("goget").Parse(`<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <meta name="go-import" content="{{.Host}}{{.ImportPath}} {{.VCS}} {{.RepoPath}}"/>
    <meta name="go-source" content="{{.SourcePath}}"/>
    <meta http-equiv="refresh" content="0; url=https://{{.Host}}{{.ImportPath}}"/>
  </head>
</html>
`))

// This is a nicely formatted css-using page for an import path.
var vanityTmpl = template.Must(template.New("vanity").Parse(`<!DOCTYPE html>
<!DOCTYPE html>
<html lang="en">
<head>
  <meta http-equiv="content-type" content="text/html; charset=UTF-8">
  <title>Package {{.Title}} - Go Lift</title>
  <link rel="icon" href="/favicon.ico" type="image/x-icon"/>

  <meta name="go-import" content="{{.Host}}{{.ImportPath}} {{.VCS}} {{.RepoPath}}"/>
  <meta name="go-source" content="{{.SourcePath}}"/>
  <meta name="description" content="{{.RepoPath}}">
  <meta name="author" content="Copyright 2019 - Go Lift">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link href="https://fonts.googleapis.com/css?family=Raleway:400,300,600" rel="stylesheet" type="text/css">

  <!-- these are in static/css -->
  <link rel="stylesheet" href="https://docs.golift.io/css/normalize.css">
  <link rel="stylesheet" href="https://docs.golift.io/css/custom.css">
  <link rel="stylesheet" href="https://docs.golift.io/css/skeleton.css">
</head>
<body>
  <div class="container">
{{if ne .ImageURL ""}}
    <!-- header image -->
    <div class="row" style="margin-top: 10%">
      <img height="150px" src="{{.ImageURL}}">
    </div>
{{end}}
    <!-- main content -->
    <div class="row" style="margin-top: 5%">
      <div class="two-thirds column">
        <h1>{{.Host}}{{.ImportPath}}</h1>
        <p>{{.Description}}</p>
      </div>
{{- if .Links}}
      <!-- custom links -->
      <div class="one-third column">
        <h4>Resources</h1>
{{- range .Links}}
        <li><a href="{{.URL}}">{{.Title}}</a></li>{{end}}
      </div>{{end}}
    </div>

    <!-- built-in links -->
    <div class="value-props row">
      <div class="one-third column value-prop">
        <form action="https://godoc.org/{{.Host}}{{.ImportPath}}" method="get">
          <input type="button" class="button button-primary" onclick="window.location.href = 'https://godoc.org/{{.Host}}{{.ImportPath}}';" value="Documentation"/>

          </input>
        </form>
      </div>
      <div class="one-third column value-prop">
        <form action="{{.RepoPath}}" method="get">
          <input type="button" class="button button-primary" onClick="window.location.href = '{{.RepoPath}}';" value="Code Repository"/>
        </form>
      </div>
    </div>

    <!-- usage & logo -->
    <div class="value-props row">
      <div class="two-thirds column value-prop">
        <p>Download this package.</p>
        <pre>
    go get {{.Host}}{{.ImportPath}}
        </pre>
        <p>Use this package.</p>
        <pre>
    import (
      "{{.Host}}{{.ImportPath}}"
    )
        </pre>
        <p>Refer to the package as <strong>{{.Title}}</strong></p>
      </div>
      <div class="one-third column value-prop">
        <img class="value-img" src="https://docs.golift.io/apple-touch-icon.png">
        <p>&copy; 2019 Go Lift<p>
      </div>
    </div>

  </div>
</body>
</html>`))
