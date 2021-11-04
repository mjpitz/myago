package internal

import (
	"text/template"

	"github.com/urfave/cli/v2"
)

const versionTemplate = "{{ .Name }} {{ range $key, $value := .Metadata }}{{ $key }}={{ $value }} {{ end }}\n"

var VersionCommand = &cli.Command{
	Name:      "version",
	Usage:     "Print the binary version information",
	UsageText: "myago version",
	Action: func(ctx *cli.Context) error {
		return template.
			Must(template.New("version").Parse(versionTemplate)).
			Execute(ctx.App.Writer, ctx.App)
	},
	HideHelpCommand: true,
}
