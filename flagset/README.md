# flagset
--
    import "github.com/mjpitz/myago/flagset"

Package flagset provides logic for creating a flagset from a struct.

## Usage

#### func  Extract

```go
func Extract(v interface{}) []cli.Flag
```
Extract parses the provided object to create a flagset.

#### func  ExtractPrefix

```go
func ExtractPrefix(prefix string, v interface{}) []cli.Flag
```
ExtractPrefix parses the provided to create a flagset with the provided Prefix.

#### type Common

```go
type Common struct {
	Name     string
	FlagName string
	Aliases  []string
	Usage    string
	EnvVars  []string
	Default  string
}
```

Common encapsulates common elements across all flag types.

#### type Extractor

```go
type Extractor struct {
	Prefix    []string
	EnvPrefix []string
}
```

Extractor extracts flags from provided interfaces.

#### func (Extractor) Child

```go
func (f Extractor) Child(name string) Extractor
```
Child creates a new Extractor and adds name to the end of the current Prefix.

#### func (Extractor) Clone

```go
func (f Extractor) Clone() Extractor
```
Clone creates a copy of the current Extractor.

#### func (Extractor) Common

```go
func (f Extractor) Common(field reflect.StructField) *Common
```
Common returns the common metadata between all fields.

#### func (Extractor) Extract

```go
func (f Extractor) Extract(v interface{}) []cli.Flag
```
Extract returns the set of flags associated with the provided structure.

#### func (Extractor) FormatBoolFlag

```go
func (f Extractor) FormatBoolFlag(common *Common, fieldValue reflect.Value) (flag *cli.BoolFlag, err error)
```
FormatBoolFlag creates a cli.BoolFlag for the given common configuration and
value.

#### func (Extractor) FormatDurationFlag

```go
func (f Extractor) FormatDurationFlag(common *Common, fieldValue reflect.Value) (flag *cli.DurationFlag, err error)
```
FormatDurationFlag creates a cli.DurationFlag for the given common configuration
and value.

#### func (Extractor) FormatFlag

```go
func (f Extractor) FormatFlag(common *Common, value reflect.Value) (flag cli.Flag, err error)
```
FormatFlag attempts to create a cli.Flag based on the type of the value.

#### func (Extractor) FormatIntFlag

```go
func (f Extractor) FormatIntFlag(common *Common, fieldValue reflect.Value) (flag *cli.IntFlag, err error)
```
FormatIntFlag creates a cli.IntFlag for the given common configuration and
value.

#### func (Extractor) FormatStringFlag

```go
func (f Extractor) FormatStringFlag(common *Common, fieldValue reflect.Value) (flag *cli.StringFlag, err error)
```
FormatStringFlag creates a cli.StringFlag for the given common configuration and
value.

#### func (Extractor) FormatStringSliceFlag

```go
func (f Extractor) FormatStringSliceFlag(common *Common, fieldValue reflect.Value) (flag *cli.StringSliceFlag, err error)
```
FormatStringSliceFlag creates a cli.StringSliceFlag for the given common
configuration and value.
