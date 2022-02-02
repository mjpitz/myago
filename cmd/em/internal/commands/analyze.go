// Copyright (C) 2022 Mya Pitzeruse
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

package commands

import (
	"fmt"

	"github.com/mjpitz/myago/cmd/em/internal/index"
	"github.com/mjpitz/myago/cmd/em/internal/jenkins"
	"github.com/mjpitz/myago/flagset"
	"github.com/urfave/cli/v2"
)

type AnalyzeConfig struct {
	Index   index.Config   `json:"index"`
	Jenkins jenkins.Config `json:"jenkins"`
}

var (
	analyzeConfig = &AnalyzeConfig{
		Jenkins: jenkins.Config{
			Jobs: cli.NewStringSlice(),
		},
	}

	Analyze = &cli.Command{
		Name:      "analyze",
		Usage:     "Generate data sets for a variety of integrations.",
		UsageText: "em analyze <integration>",
		Flags:     flagset.ExtractPrefix("em", analyzeConfig),
		Action: func(ctx *cli.Context) error {
			idx, err := index.Open(analyzeConfig.Index)
			if err != nil {
				return err
			}
			defer idx.Close()

			integration := ctx.Args().Get(0)
			switch integration {
			case "":
				return fmt.Errorf("missing integration")
			case "jenkins":
				return jenkins.Run(ctx.Context, analyzeConfig.Jenkins, idx)
			default:
				return fmt.Errorf("unknonw integration: %s", integration)
			}
		},
		HideHelpCommand: true,
	}
)
