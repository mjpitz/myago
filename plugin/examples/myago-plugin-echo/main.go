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
	"log"

	"go.pitz.tech/lib/plugin"
	"go.pitz.tech/lib/yarpc"
)

func main() {
	yarpc.HandleFunc("/echo", func(stream yarpc.Stream) (err error) {
		msg := map[string]interface{}{}
		err = stream.ReadMsg(&msg)
		if err != nil {
			return
		}

		err = stream.WriteMsg(msg)
		if err != nil {
			return
		}

		return nil
	})

	err := yarpc.Serve(plugin.Listen())
	if err != nil {
		log.Fatal(err)
	}
}
