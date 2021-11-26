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

package internal

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var phoneCodeMap = map[string]string {
	"a": "2", "b": "2", "c": "2",
	"d": "3", "e": "3", "f": "3",
	"g": "4", "h": "4", "i": "4",
	"j": "5", "k": "5", "l": "5",
	"m": "6", "n": "6", "o": "6",
	"p": "7", "q": "7", "r": "7", "s": "7",
	"t": "8", "u": "8", "v": "8",
	"w": "9", "x": "9", "y": "9", "z": "9",
}

var EncodeCommand = &cli.Command{
	Name:      "encode",
	Usage:     "Encode is a simple utility to encode a provided argument using the specified encoding.",
	UsageText: "myago encode <plaintext>",
	Action: func(ctx *cli.Context) error {
		plaintext := ctx.Args().Get(0)
		cyphertext := ""

		for _, ch := range strings.Split(plaintext, "") {
			n, ok := phoneCodeMap[ch]
			if !ok {
				return fmt.Errorf("unrecognized char: %s", ch)
			}

			cyphertext += n
		}

		_, err := ctx.App.Writer.Write([]byte(cyphertext))
		return err
	},
	HideHelpCommand: true,
}
