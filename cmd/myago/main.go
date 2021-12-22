// Copyright (C) 2021 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/mjpitz/myago/cmd/myago/internal/commands"
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
			commands.EncodeCommand,
			commands.ScaffoldCommand,
			commands.VersionCommand,
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
