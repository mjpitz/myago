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

package lazy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/lazy"
)

type Test struct{}

func TestOnce(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Load        func() (*Test, error)
		LoadContext func(ctx context.Context) (*Test, error)
		Error       string
	}{
		{
			Name:  "loader not set",
			Error: "loader not set",
		},
		{
			Name: "load failure",
			Load: func() (*Test, error) {
				return nil, fmt.Errorf("failed to load")
			},
			Error: "failed to load",
		},
		{
			Name: "load success",
			Load: func() (*Test, error) {
				return &Test{}, nil
			},
		},
		{
			Name: "load failure with context",
			LoadContext: func(ctx context.Context) (*Test, error) {
				return nil, fmt.Errorf("failed to load")
			},
			Error: "failed to load",
		},
		{
			Name: "load success with context",
			LoadContext: func(ctx context.Context) (*Test, error) {
				return &Test{}, nil
			},
		},
	}

	ctx := context.Background()
	for _, testCase := range testCases {
		t.Log(testCase.Name)

		once := &lazy.Once{}

		if testCase.LoadContext != nil {
			once.Loader = testCase.LoadContext
		} else if testCase.Load != nil {
			once.Loader = testCase.Load
		}

		v1, err := once.Get(ctx)
		if len(testCase.Error) > 0 {
			require.Error(t, err)
			require.Equal(t, testCase.Error, err.Error())
		} else {
			require.NoError(t, err)
			require.NotNil(t, v1)

			v2, err := once.Get(ctx)
			require.NoError(t, err)
			require.NotNil(t, v2)

			require.Equal(t, v1, v2)
		}
	}
}
