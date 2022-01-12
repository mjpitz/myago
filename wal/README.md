# wal

Package wal provides a simple write-ahead log implementation inspired by
Indeed's BasicRecordFile implementation. Each record in the file is stored using
the following format:

    [length - varint][record content][checksum]

Unlike the reference implementation, the record length is written as a varint to
help conserve space. The checksum is a simple CRC32 checksum. Reference:
https://github.com/indeedeng/lsmtree/blob/master/recordlog/src/main/java/com/indeed/lsmtree/recordlog/BasicRecordFile.java

```go
import github.com/mjpitz/myago/wal
```

## Usage

#### type Reader

```go
type Reader struct {
}
```

Reader implements the logic for reading information from the write-ahead log.
The underlying file is wrapped with a buffered reader to help improve
performance.

#### func OpenReader

```go
func OpenReader(ctx context.Context, filepath string) (*Reader, error)
```

OpenReader opens a new read-only handle to the target file.

#### func (\*Reader) Close

```go
func (r *Reader) Close() error
```

#### func (\*Reader) Position

```go
func (r *Reader) Position() uint64
```

Position returns the current position of the reader.

#### func (\*Reader) Read

```go
func (r *Reader) Read(p []byte) (n int, err error)
```

#### func (\*Reader) Seek

```go
func (r *Reader) Seek(offset int64, whence int) (int64, error)
```

#### type Writer

```go
type Writer struct {
}
```

Writer implements the logic for writing information to the write-ahead log. The
underlying file is wrapped with a buffered writer to help improve durability of
writes.

#### func OpenWriter

```go
func OpenWriter(ctx context.Context, filepath string) (*Writer, error)
```

OpenWriter opens a new append-only handle that writes data to the target file.

#### func (\*Writer) Close

```go
func (w *Writer) Close() error
```

#### func (\*Writer) Flush

```go
func (w *Writer) Flush() error
```

#### func (\*Writer) Sync

```go
func (w *Writer) Sync() error
```

#### func (\*Writer) Write

```go
func (w *Writer) Write(p []byte) (int, error)
```
