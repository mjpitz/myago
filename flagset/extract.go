package flagset

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func format(prefix []string, name string) string {
	val := name
	for i := 1; i <= len(prefix); i++ {
		val = prefix[len(prefix)-i] + "_" + val
	}

	return val
}

// Common encapsulates common elements across all flag types.
type Common struct {
	Name     string
	FlagName string
	Aliases  []string
	Usage    string
	EnvVars  []string
	Default  string
}

// Extractor extracts flags from provided interfaces.
type Extractor struct {
	Prefix    []string
	EnvPrefix []string
}

// Clone creates a copy of the current Extractor.
func (f Extractor) Clone() Extractor {
	nf := Extractor{}
	nf.Prefix = append(nf.Prefix, f.Prefix...)
	nf.EnvPrefix = append(nf.EnvPrefix, f.EnvPrefix...)

	return nf
}

// Child creates a new Extractor and adds name to the end of the current Prefix.
func (f Extractor) Child(name string) Extractor {
	nf := f.Clone()
	nf.Prefix = append(nf.Prefix, name)
	nf.EnvPrefix = append(nf.EnvPrefix, name)

	return nf
}

// Common returns the common metadata between all fields.
func (f Extractor) Common(field reflect.StructField) *Common {
	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "-" {
		return nil
	}

	common := &Common{
		Name:     name,
		FlagName: format(f.Prefix, name),
		Usage:    field.Tag.Get("usage"),
		EnvVars: []string{
			strings.ToUpper(format(f.EnvPrefix, name)),
		},
		Default: field.Tag.Get("default"),
	}

	if alias := field.Tag.Get("aliases"); alias != "" {
		common.Aliases = strings.Split(alias, ",")
	} else if alias := field.Tag.Get("alias"); alias != "" {
		common.Aliases = strings.Split(alias, ",")
	}

	if env := field.Tag.Get("env"); env != "" {
		common.EnvVars = append(common.EnvVars, strings.Split(env, ",")...)
	}

	return common
}

// FormatDurationFlag creates a cli.DurationFlag for the given common configuration and value.
func (f Extractor) FormatDurationFlag(common *Common, fieldValue reflect.Value) (flag *cli.DurationFlag, err error) {
	flag = &cli.DurationFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Destination: fieldValue.Addr().Interface().(*time.Duration),
	}

	if !fieldValue.IsZero() {
		flag.Value = fieldValue.Interface().(time.Duration)
	} else if common.Default != "" {
		flag.Value, err = time.ParseDuration(common.Default)
		if err != nil {
			return nil, err
		}
	}

	return flag, nil
}

// FormatStringSliceFlag creates a cli.StringSliceFlag for the given common configuration and value.
func (f Extractor) FormatStringSliceFlag(common *Common, fieldValue reflect.Value) (flag *cli.StringSliceFlag, err error) {
	flag = &cli.StringSliceFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Destination: fieldValue.Interface().(*cli.StringSlice),
	}

	if !fieldValue.IsZero() {
		flag.Value = fieldValue.Interface().(*cli.StringSlice)
	}

	return flag, nil
}

// FormatStringFlag creates a cli.StringFlag for the given common configuration and value.
func (f Extractor) FormatStringFlag(common *Common, fieldValue reflect.Value) (flag *cli.StringFlag, err error) {
	flag = &cli.StringFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Destination: fieldValue.Addr().Interface().(*string),
		Value:       common.Default,
	}

	if !fieldValue.IsZero() {
		flag.Value = fieldValue.String()
	}

	return flag, nil
}

// FormatIntFlag creates a cli.IntFlag for the given common configuration and value.
func (f Extractor) FormatIntFlag(common *Common, fieldValue reflect.Value) (flag *cli.IntFlag, err error) {
	flag = &cli.IntFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Destination: fieldValue.Addr().Interface().(*int),
	}

	if !fieldValue.IsZero() {
		flag.Value = int(fieldValue.Int())
	} else if common.Default != "" {
		flag.Value, err = strconv.Atoi(common.Default)
		if err != nil {
			return nil, err
		}
	}

	return flag, nil
}

// FormatBoolFlag creates a cli.BoolFlag for the given common configuration and value.
func (f Extractor) FormatBoolFlag(common *Common, fieldValue reflect.Value) (flag *cli.BoolFlag, err error) {
	flag = &cli.BoolFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Destination: fieldValue.Addr().Interface().(*bool),
	}

	if !fieldValue.IsZero() {
		flag.Value = fieldValue.Bool()
	} else if common.Default != "" {
		flag.Value, err = strconv.ParseBool(common.Default)
		if err != nil {
			return nil, err
		}
	}

	return flag, nil
}

// FormatFlag attempts to create a cli.Flag based on the type of the value
func (f Extractor) FormatFlag(common *Common, value reflect.Value) (flag cli.Flag, err error) {
	switch value.Interface().(type) {
	case time.Duration:
		return f.FormatDurationFlag(common, value)
	case *cli.StringSlice:
		return f.FormatStringSliceFlag(common, value)
	default:
		switch value.Type().Kind() {
		case reflect.String:
			return f.FormatStringFlag(common, value)
		case reflect.Int:
			return f.FormatIntFlag(common, value)
		case reflect.Bool:
			return f.FormatBoolFlag(common, value)
		}
	}

	return nil, nil
}

func (f Extractor) extractField(value reflect.Value, field reflect.StructField) (flags []cli.Flag) {
	common := f.Common(field)
	if common == nil {
		return
	}

	switch field.Type.Kind() {
	case reflect.Struct:
		formatter := f
		if common.Name != "" {
			formatter = f.Child(common.Name)
		}

		return append(flags, formatter.extract(value)...)
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(field.Type.Elem()))
		}
	}

	flag, err := f.FormatFlag(common, value)
	switch {
	case err != nil:
		panic(fmt.Sprintf("failed to format flag: %v", err))
	case flag != nil:
		flags = append(flags, flag)
	}

	return flags
}

func (f Extractor) extract(value reflect.Value) (flags []cli.Flag) {
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		field := value.Type().Field(i)

		flags = append(flags, f.extractField(fieldValue, field)...)
	}

	return flags
}

// Extract returns the set of flags associated with the provided structure.
func (f Extractor) Extract(v interface{}) []cli.Flag {
	return f.extract(reflect.Indirect(reflect.ValueOf(v)))
}

// Extract parses the provided object to create a flagset.
func Extract(v interface{}) []cli.Flag {
	return Extractor{}.Extract(v)
}

// ExtractPrefix parses the provided to create a flagset with the provided Prefix.
func ExtractPrefix(prefix string, v interface{}) []cli.Flag {
	return Extractor{
		EnvPrefix: []string{prefix},
	}.Extract(v)
}
