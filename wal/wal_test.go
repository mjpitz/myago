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

package wal_test

import (
	"context"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"go.pitz.tech/lib/wal"
)

func TestWAL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	vlogPath := filepath.Join(t.TempDir(), "test.vlog")

	writer, err := wal.OpenWriter(ctx, vlogPath)
	require.NoError(t, err)
	defer writer.Close()

	reader, err := wal.OpenReader(ctx, vlogPath)
	require.NoError(t, err)
	defer reader.Close()

	require.Equal(t, uint64(0), reader.Position())

	_, err = writer.Write([]byte("hello world"))
	require.NoError(t, err)

	err = writer.Sync()
	require.NoError(t, err)

	read := make([]byte, 100)
	for i := 0; i < 2; i++ {
		n, err := reader.Read(read)
		require.NoError(t, err)

		require.Equal(t, "hello world", string(read[:n]))

		require.Equal(t, uint64(0x10), reader.Position())

		pos, err := reader.Seek(0, io.SeekStart)
		require.NoError(t, err)
		require.Equal(t, int64(0), pos)
	}
}
