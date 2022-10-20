# commands




```go
import go.pitz.tech/lib/cmd/em/internal/commands
```

## Usage

```go
var (
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
```

```go
var (
	Auth = &cli.Command{
		Name:  "auth",
		Usage: "Authenticate using common mechanisms.",
		Subcommands: []*cli.Command{
			{
				Name:  "oidc",
				Usage: "Authenticate with an OIDC provider.",
				Flags: flagset.ExtractPrefix("em", oidcAuthConfig),
				Action: func(ctx *cli.Context) error {
					uri, err := url.Parse(oidcAuthConfig.RedirectURL)
					if err != nil {
						return err
					}

					svr := &http.Server{
						Addr: uri.Host,
					}

					if len(oidcAuthConfig.Scopes.Value()) == 0 {
						oidcAuthConfig.Scopes = cli.NewStringSlice("openid", "profile", "email")
					}

					cctx, cancel := context.WithCancel(ctx.Context)
					defer cancel()

					svr.Handler = oidcauth.ServeMux(*oidcAuthConfig, func(token *oauth2.Token) {
						defer cancel()

						enc := json.NewEncoder(ctx.App.Writer)
						enc.SetIndent("", "  ")
						_ = enc.Encode(token)
					})

					group := &errgroup.Group{}

					group.Go(func() error {
						time.Sleep(time.Second)
						return browser.Open(ctx.Context, uri.Scheme+"://"+uri.Host+"/login")
					})

					group.Go(svr.ListenAndServe)

					<-cctx.Done()
					err = svr.Shutdown(ctx.Context)
					_ = group.Wait()

					return nil
				},
				HideHelpCommand: true,
			},
			{
				Name:  "storj",
				Usage: "Authenticate with a Storj OIDC provider.",
				Flags: flagset.ExtractPrefix("em", storjAuthConfig),
				Action: func(ctx *cli.Context) error {
					uri, err := url.Parse(storjAuthConfig.RedirectURL)
					if err != nil {
						return err
					}

					svr := &http.Server{
						Addr: uri.Host,
						BaseContext: func(_ net.Listener) context.Context {
							return ctx.Context
						},
					}

					if len(storjAuthConfig.Scopes.Value()) == 0 {
						storjAuthConfig.Scopes = cli.NewStringSlice("openid", "profile", "email", "object:list", "object:read", "object:write", "object:delete")
					}

					cctx, cancel := context.WithCancel(ctx.Context)
					defer cancel()

					svr.Handler = storjauth.ServeMux(*storjAuthConfig, func(token *oauth2.Token, rootKey []byte) {
						defer cancel()

						enc := json.NewEncoder(ctx.App.Writer)
						enc.SetIndent("", "  ")
						_ = enc.Encode(struct {
							Token   *oauth2.Token `json:"token"`
							RootKey []byte        `json:"root_key"`
						}{
							Token:   token,
							RootKey: rootKey,
						})
					})

					group := &errgroup.Group{}

					group.Go(func() error {
						time.Sleep(time.Second)
						url := uri.Scheme + "://" + uri.Host + "/login"

						zaputil.Extract(ctx.Context).Info("Opening " + url)
						return browser.Open(ctx.Context, url)
					})

					group.Go(svr.ListenAndServe)

					<-cctx.Done()
					_ = svr.Shutdown(ctx.Context)
					_ = group.Wait()

					return nil
				},
				HideHelpCommand: true,
			},
		},
		HideHelpCommand: true,
	}
)
```

```go
var (
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
```

```go
var (
	Scaffold = &cli.Command{
		Name:  "scaffold",
		Usage: "Scaffold out a new project or add onto an existing one.",
		UsageText: flagset.ExampleString(
			"em scaffold [options] <name>",
			"em scaffold features    # will output a list of features and aliases",
			"em scaffold --mkdir --license mpl --features init <name>",
			"em scaffold --mkdir --license mpl --features init --features bin <name>",
		),
		Flags: flagset.ExtractPrefix("em", scaffoldConfig),
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return fmt.Errorf("name not specified")
			}

			name := ctx.Args().Get(0)
			if name == "features" {
				return template.Must(
					template.New("scaffold-help").
						Funcs(map[string]interface{}{
							"join": func(elems []string, sep string) string {
								return strings.Join(elems, sep)
							},
						}).
						Parse(scaffoldHelpTemplate),
				).Execute(ctx.App.Writer, map[string]interface{}{
					"features": scaffold.FilesByFeature,
					"aliases":  scaffold.FeatureAliases,
				})
			}

			if scaffoldConfig.Mkdir {
				zaputil.Extract(ctx.Context).Info("making directory")
				if err := os.MkdirAll(name, 0755); err != nil {
					return errors.Wrap(err, "failed to make project directory")
				}

				if err := os.Chdir(name); err != nil {
					return errors.Wrap(err, "failed to change into directory")
				}
			}

			zaputil.Extract(ctx.Context).Info("rendering files")
			files := scaffold.Template(
				scaffold.Data{
					Name:     name,
					License:  scaffoldConfig.License,
					Features: scaffoldConfig.Features.Value(),
				},
			).Render(ctx.Context)

			zaputil.Extract(ctx.Context).Info("writing files")
			afs := vfs.Extract(ctx.Context)
			for _, file := range files {
				dir := filepath.Dir(file.Name)
				_ = afs.MkdirAll(dir, 0755)

				if exists, _ := afero.Exists(afs, file.Name); exists {

					continue
				}

				zaputil.Extract(ctx.Context).Info("writing file", zap.String("file", file.Name))
				err := afero.WriteFile(afs, file.Name, file.Contents, 0644)
				if err != nil {
					return err
				}
			}

			if scaffoldConfig.Mkdir {
				if exists, _ := afero.Exists(afs, "go.mod"); exists {
					_, err := exec.Command("go", "mod", "tidy").CombinedOutput()
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
		HideHelpCommand: true,
	}
)
```

```go
var (
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
```

```go
var (
	ULID = &cli.Command{
		Name:  "ulid",
		Usage: "Generate or format myago/ulids.",
		Flags: flagset.ExtractPrefix("em", ulidGen),
		Subcommands: []*cli.Command{
			{
				Name:  "format",
				Usage: "Parse and format provided myago/ulids.",
				Flags: flagset.ExtractPrefix("", ulidFormat),
				Action: func(ctx *cli.Context) error {
					in, err := ioutil.ReadAll(ctx.App.Reader)
					if err != nil {
						return err
					}

					var parsed ulid.ULID

					switch ulidFormat.In {
					case "string":
						parsed, err = ulid.Parse(string(in))
					case "bytes":
						parsed = in
					default:
						err = fmt.Errorf("unrecognized input type: %s (available: string, bytes)", uuidFormat.In)
					}

					if err != nil {
						return err
					}

					switch ulidFormat.Out {
					case "json":
						enc := json.NewEncoder(ctx.App.Writer)
						enc.SetIndent("", "  ")

						err = enc.Encode(map[string]any{
							"skew":    parsed.Skew(),
							"time":    parsed.Timestamp().Local(),
							"payload": parsed.Payload(),
						})
					case "string":
						_, err = ctx.App.Writer.Write([]byte(parsed.String()))
					case "bytes":
						_, err = ctx.App.Writer.Write(parsed.Bytes())
					default:
						err = fmt.Errorf("unrecognized output type: %s (available: json, string, bytes)", uuidFormat.Out)
					}

					return err
				},
				HideHelpCommand: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			c := ctx.Context
			ulid, err := ulid.Extract(c).Generate(c, ulidGen.Size)
			if err != nil {
				return err
			}

			switch ulidGen.Out {
			case "string":
				_, err = ctx.App.Writer.Write([]byte(ulid.String()))
			case "bytes":
				_, err = ctx.App.Writer.Write(ulid.Bytes())
			default:
				err = fmt.Errorf("unrecognized output type: %s (available: string, bytes)", ulidFormat.Out)
			}

			return err
		},
		HideHelpCommand: true,
	}
)
```

```go
var Version = &cli.Command{
	Name:      "version",
	Usage:     "Print the binary version information.",
	UsageText: "em version",
	Action: func(ctx *cli.Context) error {
		return template.
			Must(template.New("version").Parse(versionTemplate)).
			Execute(ctx.App.Writer, ctx.App)
	},
	HideHelpCommand: true,
}
```

#### type AnalyzeConfig

```go
type AnalyzeConfig struct {
	Index   index.Config   `json:"index"`
	Jenkins jenkins.Config `json:"jenkins"`
}
```


#### type EncodeConfig

```go
type EncodeConfig struct {
	In  string `json:"in"  alias:"i" usage:"the input encoding"  default:"ascii"`
	Out string `json:"out" alias:"o" usage:"the output encoding" default:"ascii"`
}
```


#### type ScaffoldConfig

```go
type ScaffoldConfig struct {
	Mkdir    bool             `json:"mkdir"    usage:"specify if we should make the target project directory"`
	License  string           `json:"license"  usage:"specify which license should be applied to the project" default:"agpl3"`
	Features *cli.StringSlice `json:"features" usage:"specify the features to generate"`
}
```
