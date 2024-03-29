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

package config_test

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/config"
	"go.pitz.tech/lib/vfs"
)

const configJSON = `{
	"name": "hello",
	"description": "world"
}`

const configYAML = `
name: "hello"
description: "world"
`

const configTOML = `
name = "hello"
description = "world"
`

type TestConfig struct {
	Name        string
	Description string
}

func TestLoader(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	memFs := afero.NewMemMapFs()
	_ = afero.WriteFile(memFs, "/test/config.json", []byte(configJSON), 0o755)
	_ = afero.WriteFile(memFs, "/test/config.yaml", []byte(configYAML), 0o755)
	_ = afero.WriteFile(memFs, "/test/config.toml", []byte(configTOML), 0o755)
	_ = afero.WriteFile(memFs, "/test/unrecognized.ext", []byte(""), 0o755)
	_ = afero.WriteFile(memFs, "/test/missing-ext", []byte(""), 0o755)
	ctx = vfs.ToContext(ctx, memFs)

	tests := []string{
		"/test/config.json",
		"/test/config.yaml",
		"/test/config.toml",
	}

	for _, file := range tests {
		cfg := &TestConfig{}
		err := config.Load(ctx, cfg, file)
		require.NoError(t, err)

		require.Equal(t, "hello", cfg.Name)
		require.Equal(t, "world", cfg.Description)
	}

	negativeTests := []string{
		"/test/unrecognized.ext",
		"/test/missing-ext",
		"/test/non-existent",
	}

	for _, file := range negativeTests {
		cfg := &TestConfig{}
		err := config.Load(ctx, cfg, file)
		require.NoError(t, err)

		// no change, soft failure
		require.Equal(t, "", cfg.Name)
		require.Equal(t, "", cfg.Description)
	}
}
