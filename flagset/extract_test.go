package flagset_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/mjpitz/myago/flagset"
)

type Options struct {
	Endpoint    string        `json:"endpoint"    alias:"e" usage:"the endpoint of the server we're speaking to" default:"default-endpoint"`
	EnableSSL   bool          `json:"enable_ssl"  alias:"s" usage:"enable encryption between processes" default:"false"`
	Temperature int           `json:"temperature" alias:"t" default:"50"`
	Interval    time.Duration `json:"interval"    alias:"i" default:"5m"`
}

type Full struct {
	Options  Options          `json:"options"`
	Features *cli.StringSlice `json:"features" alias:"f"`
}

type Expectation struct {
	name  string
	alias string
	env   string
	usage string
	value interface{}
}

func verifyExpectations(t *testing.T, flags []cli.Flag, expectations []Expectation) {
	t.Helper()

	require.Len(t, flags, len(expectations))

	for i, flag := range flags {
		e := expectations[i]

		switch f := flag.(type) {
		case *cli.StringFlag:
			require.Equal(t, e.name, f.Name)
			require.Equal(t, e.alias, f.Aliases[0])
			require.Equal(t, e.env, f.EnvVars[0])
			require.Equal(t, e.usage, f.Usage)
			require.Equal(t, e.value, f.Value)
		case *cli.BoolFlag:
			require.Equal(t, e.name, f.Name)
			require.Equal(t, e.alias, f.Aliases[0])
			require.Equal(t, e.env, f.EnvVars[0])
			require.Equal(t, e.usage, f.Usage)
			require.Equal(t, e.value, f.Value)
		case *cli.IntFlag:
			require.Equal(t, e.name, f.Name)
			require.Equal(t, e.alias, f.Aliases[0])
			require.Equal(t, e.env, f.EnvVars[0])
			require.Equal(t, e.usage, f.Usage)
			require.Equal(t, e.value, f.Value)
		case *cli.DurationFlag:
			require.Equal(t, e.name, f.Name)
			require.Equal(t, e.alias, f.Aliases[0])
			require.Equal(t, e.env, f.EnvVars[0])
			require.Equal(t, e.usage, f.Usage)
			require.Equal(t, e.value, f.Value)
		}
	}
}

func TestExtract(t *testing.T) {
	t.Parallel()

	fromTag := &Full{}
	fromStruct := &Full{
		Options: Options{
			Endpoint:    "override",
			EnableSSL:   true,
			Temperature: 100,
			Interval:    10 * time.Minute,
		},
		Features: cli.NewStringSlice("awe yeah"),
	}

	testCases := []struct {
		name         string
		value        interface{}
		expectations []Expectation
	}{
		{"tag", fromTag, []Expectation{
			{"options_endpoint", "e", "OPTIONS_ENDPOINT", "the endpoint of the server we're speaking to", "default-endpoint"},
			{"options_enable_ssl", "s", "OPTIONS_ENABLE_SSL", "enable encryption between processes", false},
			{"options_temperature", "t", "OPTIONS_TEMPERATURE", "", 50},
			{"options_interval", "i", "OPTIONS_INTERVAL", "", 5 * time.Minute},
			{"features", "f", "FEATURES", "", cli.NewStringSlice()},
		}},
		{"struct", fromStruct, []Expectation{
			{"options_endpoint", "e", "OPTIONS_ENDPOINT", "the endpoint of the server we're speaking to", "override"},
			{"options_enable_ssl", "s", "OPTIONS_ENABLE_SSL", "enable encryption between processes", true},
			{"options_temperature", "t", "OPTIONS_TEMPERATURE", "", 100},
			{"options_interval", "i", "OPTIONS_INTERVAL", "", 10 * time.Minute},
			{"features", "f", "FEATURES", "", fromStruct.Features},
		}},
	}

	for _, testCase := range testCases {
		t.Log("from " + testCase.name)

		verifyExpectations(t, flagset.Extract(testCase.value), testCase.expectations)

		t.Log("from " + testCase.name + " with prefix")

		mappedExpectations := make([]Expectation, len(testCase.expectations))
		for i, expectation := range testCase.expectations {
			mappedExpectations[i] = expectation
			mappedExpectations[i].env = "PREFIX_" + mappedExpectations[i].env
		}

		verifyExpectations(t, flagset.ExtractPrefix("Prefix", testCase.value), mappedExpectations)
	}
}
