package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const libContentTemplate = `<!-- autogenerated by main.go -->
<html lang="en">
  <head>
    <meta name="go-import" content="{{ .Package }} {{ .VCS }} {{ .Repository }}">
    <meta http-equiv="refresh" content="0;URL='{{ .Repository }}'">
    <title>{{ .Package }}</title>
  </head>
  <body>
  Redirecting you to the <a href="{{ .Repository }}">project page</a>...
  </body>
</html>
`

const indexContentTemplate = `<!-- autogenerated by main.go -->
<html lang="en">
	<head>
		<title>Go Libraries</title>
	</head>
	<body>
		{{- range $generate := .Generates }}
		<div><a href="/{{ $generate.Path }}">{{ $generate.Package }}</a></div>
		{{- end }}
	</body>
</html>
`

type Generate struct {
	Path       string
	Package    string
	Repository string
	VCS        string
}

func renderFile(filename string, template *template.Template, data any) error {
	dir := filepath.Dir(filename)
	_ = os.MkdirAll(dir, 0755)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	return template.Execute(file, data)
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

func main() {
	libTemplate := must(template.New("lib").Parse(libContentTemplate))
	indexTemplate := must(template.New("index").Parse(indexContentTemplate))
	wd := must(os.Getwd())

	baseDir := flag.String("base_dir", wd, "configure the directory where this script outputs data")

	flag.Parse()

	generates := []Generate{
		{"em", "go.pitz.tech/em", "https://github.com/mjpitz/em", "git"},
		{"lib", "go.pitz.tech/lib", "https://github.com/mjpitz/myago", "git"},
		{"okit", "go.pitz.tech/okit", "https://code.pitz.tech/mya/okit", "git"},
		{"units", "go.pitz.tech/units", "https://github.com/mjpitz/units", "git"},
		{filepath.Join("gorm", "encryption"), "go.pitz.tech/gorm/encryption", "https://github.com/mjpitz/gorm-encryption", "git"},
		{"lagg", "go.pitz.tech/lagg", "https://github.com/mjpitz/lagg", "git"},
	}

	for _, generate := range generates {
		log.Println("rendering", generate.Path)

		err := renderFile(filepath.Join(*baseDir, generate.Path, "index.html"), libTemplate, generate)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	log.Println("rendering", "index.html")
	err := renderFile(filepath.Join(*baseDir, "index.html"), indexTemplate, map[string]any{"Generates": generates})
	if err != nil {
		log.Fatalln(err.Error())
	}
}
