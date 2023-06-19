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

type common struct {
	Name     string
	FlagName string
	Aliases  []string
	Usage    string
	EnvVars  []string
	Default  string
	Hidden   bool
	Required bool
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

// common returns the common metadata between all fields.
func (f Extractor) common(field reflect.StructField) *common {
	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "-" {
		return nil
	}

	hidden, _ := strconv.ParseBool(field.Tag.Get("hidden"))
	required, _ := strconv.ParseBool(field.Tag.Get("required"))

	common := &common{
		Name:     name,
		FlagName: format(f.Prefix, name),
		Usage:    field.Tag.Get("usage"),
		EnvVars: []string{
			strings.ToUpper(format(f.EnvPrefix, name)),
		},
		Default:  field.Tag.Get("default"),
		Hidden:   hidden,
		Required: required,
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

// formatDurationFlag creates a cli.DurationFlag for the given common configuration and value.
func (f Extractor) formatDurationFlag(common *common, fieldValue reflect.Value) (flag *cli.DurationFlag, err error) {
	flag = &cli.DurationFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Required:    common.Required,
		Hidden:      common.Hidden,
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

// formatStringSliceFlag creates a cli.StringSliceFlag for the given common configuration and value.
func (f Extractor) formatStringSliceFlag(common *common, fieldValue reflect.Value) (flag *cli.StringSliceFlag, err error) {
	flag = &cli.StringSliceFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Required:    common.Required,
		Hidden:      common.Hidden,
		Destination: fieldValue.Interface().(*cli.StringSlice),
	}

	if !fieldValue.IsZero() {
		flag.Value = fieldValue.Interface().(*cli.StringSlice)
	} else if common.Default != "" {
		for _, v := range strings.Split(common.Default, ",") {
			_ = flag.Value.Set(v)
		}
	}

	return flag, nil
}

// formatGenericFlag creates a cli.StringSliceFlag for the given common configuration and value.
func (f Extractor) formatGenericFlag(common *common, fieldValue reflect.Value) (flag *cli.GenericFlag, err error) {
	flag = &cli.GenericFlag{
		Name:     common.FlagName,
		Aliases:  common.Aliases,
		Usage:    common.Usage,
		EnvVars:  common.EnvVars,
		Required: common.Required,
		Hidden:   common.Hidden,
		Value:    fieldValue.Interface().(cli.Generic),
	}

	return flag, nil
}

// formatStringFlag creates a cli.StringFlag for the given common configuration and value.
func (f Extractor) formatStringFlag(common *common, fieldValue reflect.Value) (flag *cli.StringFlag, err error) {
	flag = &cli.StringFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Required:    common.Required,
		Hidden:      common.Hidden,
		Destination: fieldValue.Addr().Interface().(*string),
		Value:       common.Default,
	}

	if !fieldValue.IsZero() {
		flag.Value = fieldValue.String()
	}

	return flag, nil
}

// formatIntFlag creates a cli.IntFlag for the given common configuration and value.
func (f Extractor) formatIntFlag(common *common, fieldValue reflect.Value) (flag *cli.IntFlag, err error) {
	flag = &cli.IntFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Required:    common.Required,
		Hidden:      common.Hidden,
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

// formatBoolFlag creates a cli.BoolFlag for the given common configuration and value.
func (f Extractor) formatBoolFlag(common *common, fieldValue reflect.Value) (flag *cli.BoolFlag, err error) {
	flag = &cli.BoolFlag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Required:    common.Required,
		Hidden:      common.Hidden,
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

// formatFloatFlag creates a cli.Float64Flag for the given common configuration and value.
func (f Extractor) formatFloatFlag(common *common, fieldValue reflect.Value) (flag *cli.Float64Flag, err error) {
	flag = &cli.Float64Flag{
		Name:        common.FlagName,
		Aliases:     common.Aliases,
		Usage:       common.Usage,
		EnvVars:     common.EnvVars,
		Required:    common.Required,
		Hidden:      common.Hidden,
		Destination: fieldValue.Addr().Interface().(*float64),
	}

	if !fieldValue.IsZero() {
		flag.Value = fieldValue.Float()
	} else if common.Default != "" {
		flag.Value, err = strconv.ParseFloat(common.Default, 64)
		if err != nil {
			return nil, err
		}
	}

	return flag, nil
}

// formatFlag attempts to create a cli.Flag based on the type of the value.
func (f Extractor) formatFlag(common *common, value reflect.Value) (cli.Flag, error) {
	switch value.Interface().(type) {
	case time.Duration:
		return f.formatDurationFlag(common, value)
	case *cli.StringSlice:
		return f.formatStringSliceFlag(common, value)
	case cli.Generic:
		return f.formatGenericFlag(common, value)
	default:
		switch value.Type().Kind() {
		case reflect.String:
			return f.formatStringFlag(common, value)
		case reflect.Int:
			return f.formatIntFlag(common, value)
		case reflect.Bool:
			return f.formatBoolFlag(common, value)
		case reflect.Float64:
			return f.formatFloatFlag(common, value)
		}
	}

	return nil, nil
}

func (f Extractor) extractField(value reflect.Value, field reflect.StructField) (flags FlagSet) {
	common := f.common(field)
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

	flag, err := f.formatFlag(common, value)
	switch {
	case err != nil:
		panic(fmt.Sprintf("failed to format flag: %v", err))
	case flag != nil && common.Name != "":
		flags = append(flags, flag)
	}

	return flags
}

func (f Extractor) extract(value reflect.Value) (flags FlagSet) {
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		field := value.Type().Field(i)

		flags = append(flags, f.extractField(fieldValue, field)...)
	}

	return flags
}

// Extract returns the set of flags associated with the provided structure.
func (f Extractor) Extract(v interface{}) FlagSet {
	return f.extract(reflect.Indirect(reflect.ValueOf(v)))
}

// Extract parses the provided object to create a flagset.
func Extract(v interface{}) FlagSet {
	return Extractor{}.Extract(v)
}

// ExtractPrefix parses the provided to create a flagset with the provided environment variable prefix.
func ExtractPrefix(prefix string, v interface{}) FlagSet {
	return Extractor{
		EnvPrefix: []string{prefix},
	}.Extract(v)
}
