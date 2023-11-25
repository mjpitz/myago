// Copyright (C) 2022 Mya Pitzeruse
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

	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/flagset"
)

type test[T any] struct {
	name     string
	values   []string
	expected flagset.Slice[T]
}

func TestStringSlice(t *testing.T) {
	tests := []test[string]{
		{"empty", []string{""}, []string{}},
		{"once", []string{"v1"}, []string{"v1"}},
		{"repeated", []string{"v1", "v2"}, []string{"v1", "v2"}},
		{"json", []string{`["v1", "v2"]`}, []string{"v1", "v2"}},
		{"json + once", []string{`["v1", "v2"]`, "v3"}, []string{"v1", "v2", "v3"}},
		{"once + json", []string{"v3", `["v1", "v2"]`}, []string{"v3", "v1", "v2"}},
	}

	for _, test := range tests {
		t.Log("running " + test.name)

		ss := &flagset.Slice[string]{}
		for _, v := range test.values {
			require.NoError(t, ss.Set(v))
		}

		require.Equal(t, test.expected.String(), ss.String())
	}
}

func TestBoolSlice(t *testing.T) {
	tests := []test[bool]{
		{"empty", []string{""}, []bool{}},
		{"once", []string{"true"}, []bool{true}},
		{"repeated", []string{"true", "false"}, []bool{true, false}},
		{"json", []string{`[true, false]`}, []bool{true, false}},
		{"json + once", []string{`[true, false]`, "true"}, []bool{true, false, true}},
		{"once + json", []string{"true", `[true, false]`}, []bool{true, true, false}},
	}

	for _, test := range tests {
		t.Log("running " + test.name)

		ss := &flagset.Slice[bool]{}
		for _, v := range test.values {
			require.NoError(t, ss.Set(v))
		}

		require.Equal(t, test.expected.String(), ss.String())
	}
}

func TestIntSlice(t *testing.T) {
	tests := []test[int]{
		{"empty", []string{""}, []int{}},
		{"once", []string{"1"}, []int{1}},
		{"repeated", []string{"1", "2"}, []int{1, 2}},
		{"json", []string{`[1, 2]`}, []int{1, 2}},
		{"json + once", []string{`[1, 2]`, "3"}, []int{1, 2, 3}},
		{"once + json", []string{"3", `[1, 2]`}, []int{3, 1, 2}},
	}

	for _, test := range tests {
		t.Log("running " + test.name)

		ss := &flagset.Slice[int]{}
		for _, v := range test.values {
			require.NoError(t, ss.Set(v))
		}

		require.Equal(t, test.expected.String(), ss.String())
	}

	{
		tests := []test[int8]{
			{"int8 - empty", []string{""}, []int8{}},
			{"int8 - once", []string{"1"}, []int8{1}},
			{"int8 - repeated", []string{"1", "2"}, []int8{1, 2}},
			{"int8 - json", []string{`[1, 2]`}, []int8{1, 2}},
			{"int8 - json + once", []string{`[1, 2]`, "3"}, []int8{1, 2, 3}},
			{"int8 - once + json", []string{"3", `[1, 2]`}, []int8{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[int8]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}

	{
		tests := []test[int16]{
			{"int16 - empty", []string{""}, []int16{}},
			{"int16 - once", []string{"1"}, []int16{1}},
			{"int16 - repeated", []string{"1", "2"}, []int16{1, 2}},
			{"int16 - json", []string{`[1, 2]`}, []int16{1, 2}},
			{"int16 - json + once", []string{`[1, 2]`, "3"}, []int16{1, 2, 3}},
			{"int16 - once + json", []string{"3", `[1, 2]`}, []int16{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[int16]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}

	{
		tests := []test[int32]{
			{"int32 - empty", []string{""}, []int32{}},
			{"int32 - once", []string{"1"}, []int32{1}},
			{"int32 - repeated", []string{"1", "2"}, []int32{1, 2}},
			{"int32 - json", []string{`[1, 2]`}, []int32{1, 2}},
			{"int32 - json + once", []string{`[1, 2]`, "3"}, []int32{1, 2, 3}},
			{"int32 - once + json", []string{"3", `[1, 2]`}, []int32{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[int32]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}

	{
		tests := []test[int64]{
			{"int64 - empty", []string{""}, []int64{}},
			{"int64 - once", []string{"1"}, []int64{1}},
			{"int64 - repeated", []string{"1", "2"}, []int64{1, 2}},
			{"int64 - json", []string{`[1, 2]`}, []int64{1, 2}},
			{"int64 - json + once", []string{`[1, 2]`, "3"}, []int64{1, 2, 3}},
			{"int64 - once + json", []string{"3", `[1, 2]`}, []int64{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[int64]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}
}

func TestUintSlice(t *testing.T) {
	{
		tests := []test[uint]{
			{"uint - empty", []string{""}, []uint{}},
			{"uint - once", []string{"1"}, []uint{1}},
			{"uint - repeated", []string{"1", "2"}, []uint{1, 2}},
			{"uint - json", []string{`[1, 2]`}, []uint{1, 2}},
			{"uint - json + once", []string{`[1, 2]`, "3"}, []uint{1, 2, 3}},
			{"uint - once + json", []string{"3", `[1, 2]`}, []uint{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[uint]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}

	{
		tests := []test[uint8]{
			{"uint8 - empty", []string{""}, []uint8{}},
			{"uint8 - once", []string{"1"}, []uint8{1}},
			{"uint8 - repeated", []string{"1", "2"}, []uint8{1, 2}},
			{"uint8 - json", []string{`[1, 2]`}, []uint8{1, 2}},
			{"uint8 - json + once", []string{`[1, 2]`, "3"}, []uint8{1, 2, 3}},
			{"uint8 - once + json", []string{"3", `[1, 2]`}, []uint8{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[uint8]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}

	{
		tests := []test[uint16]{
			{"uint16 - empty", []string{""}, []uint16{}},
			{"uint16 - once", []string{"1"}, []uint16{1}},
			{"uint16 - repeated", []string{"1", "2"}, []uint16{1, 2}},
			{"uint16 - json", []string{`[1, 2]`}, []uint16{1, 2}},
			{"uint16 - json + once", []string{`[1, 2]`, "3"}, []uint16{1, 2, 3}},
			{"uint16 - once + json", []string{"3", `[1, 2]`}, []uint16{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[uint16]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}

	{
		tests := []test[uint32]{
			{"uint32 - empty", []string{""}, []uint32{}},
			{"uint32 - once", []string{"1"}, []uint32{1}},
			{"uint32 - repeated", []string{"1", "2"}, []uint32{1, 2}},
			{"uint32 - json", []string{`[1, 2]`}, []uint32{1, 2}},
			{"uint32 - json + once", []string{`[1, 2]`, "3"}, []uint32{1, 2, 3}},
			{"uint32 - once + json", []string{"3", `[1, 2]`}, []uint32{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[uint32]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}

	{
		tests := []test[uint64]{
			{"uint64 - empty", []string{""}, []uint64{}},
			{"uint64 - once", []string{"1"}, []uint64{1}},
			{"uint64 - repeated", []string{"1", "2"}, []uint64{1, 2}},
			{"uint64 - json", []string{`[1, 2]`}, []uint64{1, 2}},
			{"uint64 - json + once", []string{`[1, 2]`, "3"}, []uint64{1, 2, 3}},
			{"uint64 - once + json", []string{"3", `[1, 2]`}, []uint64{3, 1, 2}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[uint64]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}
}

func TestFloatSlice(t *testing.T) {
	{
		tests := []test[float32]{
			{"float32 - empty", []string{""}, []float32{}},
			{"float32 - once", []string{"1.26"}, []float32{1.26}},
			{"float32 - repeated", []string{"1.26", "2.14"}, []float32{1.26, 2.14}},
			{"float32 - json", []string{`[1.26, 2.14]`}, []float32{1.26, 2.14}},
			{"float32 - json + once", []string{`[1.26, 2.14]`, "3.89"}, []float32{1.26, 2.14, 3.89}},
			{"float32 - once + json", []string{"3.89", `[1.26, 2.14]`}, []float32{3.89, 1.26, 2.14}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[float32]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}

	{
		tests := []test[float64]{
			{"float64 - empty", []string{""}, []float64{}},
			{"float64 - once", []string{"1.26"}, []float64{1.26}},
			{"float64 - repeated", []string{"1.26", "2.14"}, []float64{1.26, 2.14}},
			{"float64 - json", []string{`[1.26, 2.14]`}, []float64{1.26, 2.14}},
			{"float64 - json + once", []string{`[1.26, 2.14]`, "3.89"}, []float64{1.26, 2.14, 3.89}},
			{"float64 - once + json", []string{"3.89", `[1.26, 2.14]`}, []float64{3.89, 1.26, 2.14}},
		}

		for _, test := range tests {
			t.Log("running " + test.name)

			ss := &flagset.Slice[float64]{}
			for _, v := range test.values {
				require.NoError(t, ss.Set(v))
			}

			require.Equal(t, test.expected.String(), ss.String())
		}
	}
}
