package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/mjpitz/myago/cmd/myago/internal"
	"github.com/mjpitz/myago/flagset"
	"github.com/mjpitz/myago/zaputil"
)

type GlobalConfig struct {
	Log zaputil.Config `json:"log"`
}

func main() {
	config := &GlobalConfig{
		Log: zaputil.DefaultConfig(),
	}

	app := &cli.App{
		Name:      "myago",
		UsageText: "myago [options] <command>",
		Flags:     flagset.Extract(config),
		Commands: []*cli.Command{
			internal.ScaffoldCommand,
			internal.VersionCommand,
		},
		Before: func(ctx *cli.Context) error {
			ctx.Context = zaputil.Setup(ctx.Context, config.Log)
			return nil
		},
		HideVersion:          true,
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		BashComplete:         cli.DefaultAppComplete,
		Metadata: map[string]interface{}{
			"arch":       runtime.GOARCH,
			"go_version": strings.TrimPrefix(runtime.Version(), "go"),
			"os":         runtime.GOOS,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
