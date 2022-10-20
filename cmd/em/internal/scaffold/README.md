# scaffold




```go
import go.pitz.tech/lib/cmd/em/internal/scaffold
```

## Usage

```go
var (

	// FeatureAliases provides common name associates to their appropriate feature name.
	FeatureAliases = map[string][]string{
		"init": {"make", "git", "legal", "authors"},
		"bin":  {"gobinary", "go", "version"},
	}

	// FilesByFeature contains a map of features to the associated files they render.
	FilesByFeature = map[string][]string{
		"Authors": {
			"AUTHORS.tmpl",
		},
		"Docker": {
			"cmd/{{ .Name }}/docker-compose.yaml.tmpl",
			"cmd/{{ .Name }}/Dockerfile.tmpl",
		},
		"GoBinary": {
			"internal/commands/version.go.tmpl",
			"cmd/{{ .Name }}/main.go.tmpl",
		},
		"GoReleaser": {
			"scripts/dist-go.sh.tmpl",
			".goreleaser.yaml.tmpl",
		},
		"Go": {
			"go.mod.tmpl",
		},
		"Version": {
			"package.json.tmpl",
		},
		"Legal": {
			"legal/header.txt.tmpl",
			"LICENSE.tmpl",
		},
		"Git": {
			".gitignore.tmpl",
			".git/HEAD",
			".git/config",
		},
		"Make": {
			"Makefile.tmpl",
		},
	}
)
```

#### type Data

```go
type Data struct {
	Name     string   `json:"name"`
	License  string   `json:"license"`
	Features []string `json:"features"`
}
```

Data defines the information needed to render the template.

#### type File

```go
type File struct {
	Name     string
	Contents []byte
}
```

File defines a generic file to render.

#### type Renderer

```go
type Renderer struct {
}
```

Renderer encapsulates the business logic for rendering template data.

#### func  Template

```go
func Template(data Data) *Renderer
```
Template returns a new Renderer instance that can be used to render template
data based on specified features.

#### func (*Renderer) Render

```go
func (s *Renderer) Render(ctx context.Context) []File
```
