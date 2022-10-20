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

/*
Package yarpc implements "yet another RPC framework" on top of HashiCorp's yamux library. I wanted something with the
simplicity of Go's HTTP library and the ability to easily manage connections like gRPC.

Why? gRPC comes with a rather large foot-print and in many of these cases, I wanted a slimmer package for passing
messages between processes.

Example Server:

	type Stat struct {
		Name string
		Value int
	}

	start := time.Now()

	yarpc.HandleFunc("admin.stats", func(stream yarpc.Stream) error {
		for {
			err = stream.SendMsg(&Stat{ "uptime", time.Since(start).Seconds() })
			if err != nil {
				return err
			}
			time.Sleep(5 * time.Second)
		}
	})

	yarpc.ListenAndServe("tcp", "0.0.0.0:8080")

Example ClientConn:

	ctx := context.Background()
	conn := yarpc.Dial("tcp", "localhost:8080")

	stream := conn.openStream(ctx, "admin.stats")

	stat := Stat{}
	for {
		err = stream.RecvMsg(&stat)
		if err != nil {
			break
		}

		stat.Name // "uptime"
		stat.Name // "uptime"
	}
*/
package yarpc
