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

	return &Reader{handle, bufio.NewReader(handle)}, nil
}

// Reader implements the logic for reading information from the write-ahead log. The underlying file is wrapped with a
// buffered reader to help improve performance.
type Reader struct {
	handle afero.File
	buffer *bufio.Reader
}

func (r *Reader) Read(p []byte) (n int, err error) {
	length, err := binary.ReadUvarint(r.buffer)
	if err != nil {
		return 0, fmt.Errorf("bad record value")
	}

	buf := make([]byte, length)
	n, err = r.buffer.Read(buf)
	if err != nil {
		return 0, err
	}

	buf = buf[:n]

	var checksum uint32
	err = binary.Read(r.buffer, binary.BigEndian, &checksum)
	if err != nil {
		return 0, err
	}

	if crc32.ChecksumIEEE(buf) != checksum {
		return 0, fmt.Errorf("corrupted block")
	}

	return copy(p, buf), nil
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	pos, err := r.handle.Seek(offset, whence)
	if err != nil {
		return pos, err
	}

	r.buffer.Reset(r.handle)
	return pos, nil
}

func (r *Reader) Close() error {
	return r.handle.Close()
}

var _ io.ReadSeekCloser = &Reader{}
