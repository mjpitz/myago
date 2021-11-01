package flagset

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func format(prefix []string, name string) string {
	envVar := name
	for i := 1; i <= len(prefix); i++ {
		envVar = prefix[len(prefix)-i] + "_" + envVar
	}
	return envVar
}

func extract(prefix, envPrefix []string, value reflect.Value) []cli.Flag {
	flags := make([]cli.Flag, 0)

	for i := 0; i < value.NumField(); i++ {
		fieldValue := reflect.Indirect(value.Field(i))
		field := value.Type().Field(i)

		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" {
			continue
		}

		// all other data types
		var aliases []string
		if alias := field.Tag.Get("aliases"); alias != "" {
			aliases = strings.Split(alias, ",")
		}

		flagName := format(prefix, name)
		defaultTag := field.Tag.Get("default")

		var err error
		switch field.Type.Kind() {
		case reflect.Struct, reflect.Ptr:
			pre := append([]string{}, prefix...)
			envPre := append([]string{}, envPrefix...)

			if name != "" {
				pre = append(pre, name)
				envPre = append(envPre, name)
			}

			flags = append(flags, extract(pre, envPre, fieldValue)...)
		case reflect.String:
			flag := &cli.StringFlag{
				Name:        flagName,
				Aliases:     aliases,
				Usage:       field.Tag.Get("usage"),
				EnvVars:     []string{strings.ToUpper(flagName)},
				Destination: fieldValue.Addr().Interface().(*string),
				Value:       defaultTag,
			}

			if !fieldValue.IsZero() {
				flag.Value = fieldValue.String()
			}

			flags = append(flags, flag)
		case reflect.Int:
			flag := &cli.IntFlag{
				Name:        flagName,
				Aliases:     aliases,
				Usage:       field.Tag.Get("usage"),
				EnvVars:     []string{strings.ToUpper(flagName)},
				Destination: fieldValue.Addr().Interface().(*int),
			}

			if !fieldValue.IsZero() {
				flag.Value = int(fieldValue.Int())
			} else if defaultTag != "" {
				flag.Value, err = strconv.Atoi(defaultTag)
				if err != nil {
					panic(fmt.Sprintf("invalid int default: %s", defaultTag))
				}
			}

			flags = append(flags, flag)
		case reflect.Bool:
			flag := &cli.BoolFlag{
				Name:        flagName,
				Aliases:     aliases,
				Usage:       field.Tag.Get("usage"),
				EnvVars:     []string{strings.ToUpper(flagName)},
				Destination: fieldValue.Addr().Interface().(*bool),
			}

			if !fieldValue.IsZero() {
				flag.Value = fieldValue.Bool()
			} else if defaultTag != "" {
				flag.Value, err = strconv.ParseBool(defaultTag)
				if err != nil {
					panic(fmt.Sprintf("invalid bool default: %s", defaultTag))
				}
			}

			flags = append(flags, flag)
		}
	}

	return flags
}

// Extract parses the provided object to create a flagset.
func Extract(v interface{}) []cli.Flag {
	return extract(nil, nil, reflect.Indirect(reflect.ValueOf(v)))
}

// ExtractPrefix parses the provided to create a flagset with the provided prefix.
func ExtractPrefix(prefix string, v interface{}) []cli.Flag {
	return extract(nil, []string{prefix}, reflect.Indirect(reflect.ValueOf(v)))
}
