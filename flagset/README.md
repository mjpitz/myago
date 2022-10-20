# flagset

Package flagset provides an opinionated approach to constructing an
applications' configuration using Golang structs and tags. It's designed in a
way that allows configuration to be loaded from files, environment variables,
and/or command line flags. The following details the various tags that can be
specified on a primitive field.

- `json` - `string` - Configure the name of the flag. Convention is to use snake
  case.

- `usage` - `string` - Configure the description string of the flag.

- `default` - `any` - Configure the default value for the flag. Can be
  overridden by setting the value on the struct.

- `hidden` - `bool` - Hides the flag from output. The value can still be
  configured.

- `required` - `bool` - Specifies that the flag must be specified.

Nested structures are supported, making application configuration composable and
portable between systems.

```go
import go.pitz.tech/lib/flagset
```

## Usage

#### func ExampleString

```go
func ExampleString(examples ...string) string
```

ExampleString formats a list of examples so that they display properly in the
terminal. This function just pulls things out into a simple helper.

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

#### func (Extractor) Extract

```go
func (f Extractor) Extract(v interface{}) FlagSet
```

Extract returns the set of flags associated with the provided structure.

#### type Filter

```go
type Filter func(flag cli.Flag) bool
```

Filter allows the user to inspect the flag to determine if it should be in the
resulting FlagSet.

#### type FlagSet

```go
type FlagSet []cli.Flag
```

FlagSet provides additional functionality on top of a collection of flags.

#### func Extract

```go
func Extract(v interface{}) FlagSet
```

Extract parses the provided object to create a flagset.

#### func ExtractPrefix

```go
func ExtractPrefix(prefix string, v interface{}) FlagSet
```

ExtractPrefix parses the provided to create a flagset with the provided
environment variable prefix.

#### func (FlagSet) Filter

```go
func (flags FlagSet) Filter(allow Filter) FlagSet
```

Filter returns a new FlagSet that contains flags allowed by the provided Filter.
