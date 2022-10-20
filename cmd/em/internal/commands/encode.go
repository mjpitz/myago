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

package commands

import (
	"bufio"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"go.pitz.tech/lib/cmd/em/internal/phone"
	"go.pitz.tech/lib/flagset"
)

type EncodeConfig struct {
	In  string `json:"in"  alias:"i" usage:"the input encoding"  default:"ascii"`
	Out string `json:"out" alias:"o" usage:"the output encoding" default:"ascii"`
}

var (
	encodeConfig = &EncodeConfig{}

	Encode = &cli.Command{
		Name:      "encode",
		Usage:     "Read and write different encodings.",
		UsageText: "em encode [message]",
		Flags:     flagset.ExtractPrefix("em", encodeConfig),
		Aliases:   []string{"enc"},
		Action: func(ctx *cli.Context) error {
			writer := bufio.NewWriter(ctx.App.Writer)

			var reader io.Reader = bufio.NewReader(os.Stdin)
			if ctx.NArg() > 0 {
				reader = strings.NewReader(ctx.Args().Get(0))
			}

			decoder := reader
			switch encodeConfig.In {
			case "base64", "b64":
				decoder = base64.NewDecoder(base64.StdEncoding, reader)
			case "base64url", "b64url":
				decoder = base64.NewDecoder(base64.URLEncoding, reader)
			case "base32", "b32":
				decoder = base32.NewDecoder(base32.StdEncoding, reader)
			case "base32hex", "b32hex":
				decoder = base32.NewDecoder(base32.HexEncoding, reader)
			case "hex":
				decoder = hex.NewDecoder(reader)
			}

			var encoder io.Writer = writer
			switch encodeConfig.Out {
			case "base64", "b64":
				encoder = base64.NewEncoder(base64.StdEncoding, writer)
			case "base64url", "b64url":
				encoder = base64.NewEncoder(base64.URLEncoding, writer)
			case "base32", "b32":
				encoder = base32.NewEncoder(base32.StdEncoding, writer)
			case "base32hex", "b32hex":
				encoder = base32.NewEncoder(base32.HexEncoding, writer)
			case "hex":
				encoder = hex.NewEncoder(writer)
			case "phone":
				encoder = phone.NewEncoder(writer)
			}

			defer func() {
				defer writer.Flush()

				if readCloser, rcOK := decoder.(io.Closer); rcOK {
					_ = readCloser.Close()
				}

				if writeCloser, wcOK := encoder.(io.Closer); wcOK {
					_ = writeCloser.Close()
				}
			}()

			_, err := io.Copy(encoder, decoder)
			switch {
			case err == io.EOF:
			case err != nil:
				return err
			}

			return nil
		},
		HideHelpCommand: true,
	}
)
