package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const inlineTemplate = `<!-- autogenerated by main.go -->
<html lang="en">
  <head>
    <meta name="go-import" content="{{ .Package }} {{ .VCS }} {{ .Repository }}">
    <meta http-equiv="refresh" content="0;URL='{{ .Repository }}'">
    <title></title>
  </head>
  <body>
  Redirecting you to the <a href="{{ .Repository }}">project page</a>...
  </body>
</html>
`

type Generate struct {
	Path       string
	Package    string
	Repository string
	VCS        string
}

func renderFile(filename string, template *template.Template, generate Generate) error {
	dir := filepath.Dir(filename)
	_ = os.MkdirAll(dir, 0755)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	return template.Execute(file, generate)
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

func main() {
	template := must(template.New("inline").Parse(inlineTemplate))
	wd := must(os.Getwd())

	baseDir := flag.String("base_dir", wd, "configure the directory where this script outputs data")

	flag.Parse()

	generates := []Generate{
		{"em", "go.pitz.tech/em", "https://github.com/mjpitz/em", "git"},
		{"lib", "go.pitz.tech/lib", "https://github.com/mjpitz/myago", "git"},
	}

	for _, generate := range generates {
		log.Println("rendering", generate.Path)

		err := renderFile(filepath.Join(*baseDir, generate.Path, "index.html"), template, generate)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
