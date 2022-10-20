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
	"io/ioutil"

	"github.com/urfave/cli/v2"

	"go.pitz.tech/lib/flagset"

	"storj.io/common/uuid"
)

type uuidGenConfig struct {
	Out string `json:"out" alias:"o" usage:"specify the output format (string or bytes)"`
}

type uuidFormatConfig struct {
	uuidGenConfig

	In string `json:"in" alias:"i" usage:"specify the input format (string or bytes)"`
}

var (
	uuidGen = &uuidGenConfig{
		Out: "string",
	}

	uuidFormat = &uuidFormatConfig{
		In: "string",
		uuidGenConfig: uuidGenConfig{
			Out: "bytes",
		},
	}

	Storj = &cli.Command{
		Name:  "storj",
		Usage: "Utility scripts for working with storj-specific semantics.",
		Subcommands: []*cli.Command{
			{
				Name:  "uuid",
				Usage: "Format storj-specific UUID.",
				Flags: flagset.ExtractPrefix("em", uuidGen),
				Subcommands: []*cli.Command{
					{
						Name:  "format",
						Usage: "Swap between different formats of the UUID (string and bytes)",
						Flags: flagset.ExtractPrefix("em", uuidFormat),
						Action: func(ctx *cli.Context) error {
							in, err := ioutil.ReadAll(ctx.App.Reader)
							if err != nil {
								return err
							}

							var parsed uuid.UUID

							switch uuidFormat.In {
							case "string":
								parsed, err = uuid.FromString(string(in))
							case "bytes":
								parsed, err = uuid.FromBytes(in)
							default:
								err = fmt.Errorf("unrecognized input type: %s (available: string, bytes)", uuidFormat.In)
							}

							if err != nil {
								return err
							}

							switch uuidFormat.Out {
							case "string":
								_, err = ctx.App.Writer.Write([]byte(parsed.String()))
							case "bytes":
								_, err = ctx.App.Writer.Write(parsed.Bytes())
							default:
								err = fmt.Errorf("unrecognized output type: %s (available: string, bytes)", uuidFormat.Out)
							}

							return err
						},
						HideHelpCommand: true,
					},
				},
				Action: func(ctx *cli.Context) error {
					uuid, err := uuid.New()
					if err != nil {
						return err
					}

					switch uuidGen.Out {
					case "string":
						_, err = ctx.App.Writer.Write([]byte(uuid.String()))
					case "bytes":
						_, err = ctx.App.Writer.Write(uuid.Bytes())
					default:
						err = fmt.Errorf("unrecognized output type: %s (available: string, bytes)", uuidFormat.Out)
					}

					return err
				},
				HideHelpCommand: true,
			},
		},
		HideHelpCommand: true,
	}
)
