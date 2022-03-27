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

package wal

import (
	"bufio"
	"context"
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"

	"github.com/spf13/afero"

	"github.com/mjpitz/myago/vfs"
)

// OpenWriter opens a new append-only handle that writes data to the target file.
func OpenWriter(ctx context.Context, filepath string) (*Writer, error) {
	afs := vfs.Extract(ctx)

	handle, err := afs.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &Writer{handle, bufio.NewWriter(handle)}, nil
}

// Writer implements the logic for writing information to the write-ahead log. The underlying file is wrapped with a
// buffered writer to help improve durability of writes.
type Writer struct {
	handle afero.File
	buffer *bufio.Writer
}

func (w *Writer) Write(p []byte) (int, error) {
	length := len(p)
	checksum := crc32.ChecksumIEEE(p)

	buffer := make([]byte, 10+length+4)
	n := binary.PutUvarint(buffer, uint64(length))
	copy(buffer[n:], p)
	binary.BigEndian.PutUint32(buffer[n+length:], checksum)

	_, err := w.buffer.Write(buffer[:n+length+4])
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (w *Writer) Flush() error {
	return w.buffer.Flush()
}

func (w *Writer) Sync() error {
	return w.buffer.Flush()
}

func (w *Writer) Close() error {
	w.Flush()

	return w.handle.Close()
}

var _ io.WriteCloser = &Writer{}
