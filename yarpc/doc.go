/*
Package yarpc implements "yet another RPC framework" on top of HashiCorp's yamux library. I wanted something with the
simplicity of Go's HTTP library and the ability to easily manage connections like gRPC.

Why? gRPC comes with a rather large foot print and in many of these cases, I wanted a slimmer package for passing
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
