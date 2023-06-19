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

package flagset_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"go.pitz.tech/lib/flagset"
)

type Nested struct {
	Tagged   string `json:"tagged" hidden:"true" required:"true"`
	Untagged string
}

type Options struct {
	Endpoint    string        `json:"endpoint"    alias:"e" usage:"the endpoint of the server we're speaking to" default:"default-endpoint"`
	EnableSSL   bool          `json:"enable_ssl"  alias:"s" usage:"enable encryption between processes" default:"false"`
	Temperature int           `json:"temperature" alias:"t" default:"50"`
	Interval    time.Duration `json:"interval"    alias:"i" default:"5m"`
	Percentage  float64       `json:"percentage" default:"5.0"`
}

type Full struct {
	Options     Options          `json:"options"`
	Features    *cli.StringSlice `json:"features" alias:"f"`
	Percentiles *cli.IntSlice    `json:"percentiles"`
	Nested
}

type Expectation struct {
	name     string
	alias    string
	env      string
	usage    string
	value    interface{}
	hidden   bool
	required bool
}

func verifyExpectations(t *testing.T, flags []cli.Flag, expectations []Expectation) {
	t.Helper()

	require.Len(t, flags, len(expectations))

	for i, flag := range flags {
		e := expectations[i]

		switch f := flag.(type) {
		case *cli.StringFlag:
			require.Equal(t, e.name, f.Name)
			require.Equal(t, e.env, f.EnvVars[0])
			require.Equal(t, e.usage, f.Usage)
			require.Equal(t, e.value, f.Value)
			require.Equal(t, e.hidden, f.Hidden)
			require.Equal(t, e.required, f.Required)

			if e.alias != "" {
				require.Equal(t, e.alias, f.Aliases[0])
			}
		case *cli.BoolFlag:
			require.Equal(t, e.name, f.Name)
			require.Equal(t, e.env, f.EnvVars[0])
			require.Equal(t, e.usage, f.Usage)
			require.Equal(t, e.value, f.Value)
			require.Equal(t, e.hidden, f.Hidden)
			require.Equal(t, e.required, f.Required)

			if e.alias != "" {
				require.Equal(t, e.alias, f.Aliases[0])
			}
		case *cli.IntFlag:
			require.Equal(t, e.name, f.Name)
			require.Equal(t, e.env, f.EnvVars[0])
			require.Equal(t, e.usage, f.Usage)
			require.Equal(t, e.value, f.Value)
			require.Equal(t, e.hidden, f.Hidden)
			require.Equal(t, e.required, f.Required)

			if e.alias != "" {
				require.Equal(t, e.alias, f.Aliases[0])
			}
		case *cli.DurationFlag:
			require.Equal(t, e.name, f.Name)
			require.Equal(t, e.env, f.EnvVars[0])
			require.Equal(t, e.usage, f.Usage)
			require.Equal(t, e.value, f.Value)
			require.Equal(t, e.hidden, f.Hidden)
			require.Equal(t, e.required, f.Required)

			if e.alias != "" {
				require.Equal(t, e.alias, f.Aliases[0])
			}
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
			Percentage:  10.0,
		},
		Features:    cli.NewStringSlice("awe yeah"),
		Percentiles: cli.NewIntSlice(75, 90, 95, 97, 99),
	}

	testCases := []struct {
		name         string
		value        interface{}
		expectations []Expectation
	}{
		{"tag", fromTag, []Expectation{
			{"options_endpoint", "e", "OPTIONS_ENDPOINT", "the endpoint of the server we're speaking to", "default-endpoint", false, false},
			{"options_enable_ssl", "s", "OPTIONS_ENABLE_SSL", "enable encryption between processes", false, false, false},
			{"options_temperature", "t", "OPTIONS_TEMPERATURE", "", 50, false, false},
			{"options_interval", "i", "OPTIONS_INTERVAL", "", 5 * time.Minute, false, false},
			{"options_percentage", "", "OPTIONS_PERCENTAGE", "", 5.0, false, false},
			{"features", "f", "FEATURES", "", cli.NewStringSlice(), false, false},
			{"percentiles", "", "PERCENTILES", "", cli.NewIntSlice(), false, false},
			{"tagged", "", "TAGGED", "", "", true, true},
		}},
		{"struct", fromStruct, []Expectation{
			{"options_endpoint", "e", "OPTIONS_ENDPOINT", "the endpoint of the server we're speaking to", "override", false, false},
			{"options_enable_ssl", "s", "OPTIONS_ENABLE_SSL", "enable encryption between processes", true, false, false},
			{"options_temperature", "t", "OPTIONS_TEMPERATURE", "", 100, false, false},
			{"options_interval", "i", "OPTIONS_INTERVAL", "", 10 * time.Minute, false, false},
			{"options_percentage", "", "OPTIONS_PERCENTAGE", "", 10.0, false, false},
			{"features", "f", "FEATURES", "", fromStruct.Features, false, false},
			{"percentiles", "", "PERCENTILES", "", fromStruct.Percentiles, false, false},
			{"tagged", "", "TAGGED", "", "", true, true},
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
