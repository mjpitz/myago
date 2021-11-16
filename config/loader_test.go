package config_test

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/mjpitz/myago/config"
	"github.com/mjpitz/myago/vfs"
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
	_ = afero.WriteFile(memFs, "/test/config.json", []byte(configJSON), 0755)
	_ = afero.WriteFile(memFs, "/test/config.yaml", []byte(configYAML), 0755)
	_ = afero.WriteFile(memFs, "/test/config.toml", []byte(configTOML), 0755)
	_ = afero.WriteFile(memFs, "/test/unrecognized.ext", []byte(""), 0755)
	_ = afero.WriteFile(memFs, "/test/missing-ext", []byte(""), 0755)
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
