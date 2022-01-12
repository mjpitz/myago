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
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"

	"github.com/spf13/afero"

	"github.com/mjpitz/myago/vfs"
)

// OpenReader opens a new read-only handle to the target file.
func OpenReader(ctx context.Context, filepath string) (*Reader, error) {
	afs := vfs.Extract(ctx)

	handle, err := afs.Open(filepath)
	if err != nil {
		return nil, err
	}

	return &Reader{handle, bufio.NewReader(handle), 0}, nil
}

// Reader implements the logic for reading information from the write-ahead log. The underlying file is wrapped with a
// buffered reader to help improve performance.
type Reader struct {
	handle   afero.File
	buffer   *bufio.Reader
	position uint64
}

// Position returns the current position of the reader.
func (r *Reader) Position() uint64 {
	return r.position
}

func (r *Reader) Read(p []byte) (n int, err error) {
	data, err := r.buffer.Peek(10)
	if err != nil {
		return 0, err
	}

	buffer := bytes.NewReader(data)
	length, err := binary.ReadUvarint(buffer)
	if err != nil {
		return 0, err
	}

	n = len(data) - buffer.Len()
	data = make([]byte, uint64(n)+length+4)

	_, err = r.buffer.Read(data)
	if err != nil {
		return 0, err
	}
	r.position += uint64(n) + length + 4

	record := data[n : uint64(n)+length]
	checksum := binary.BigEndian.Uint32(data[uint64(n)+length:])

	if crc32.ChecksumIEEE(record) != checksum {
		return 0, fmt.Errorf("corrupted block")
	}

	return copy(p, record), nil
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	pos, err := r.handle.Seek(offset, whence)
	if err != nil {
		return pos, err
	}

	r.buffer.Reset(r.handle)
	r.position = uint64(pos)
	return pos, nil
}

func (r *Reader) Close() error {
	return r.handle.Close()
}

var _ io.ReadSeekCloser = &Reader{}
