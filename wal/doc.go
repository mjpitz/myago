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

// Package wal provides a simple write-ahead log implementation inspired by Indeed's BasicRecordFile implementation.
// Each record in the file is stored using the following format:
//
//    [length - varint][record content][checksum]
//
// Unlike the reference implementation, the record length is written as a varint to help conserve space. The checksum is
// a simple CRC32 checksum. Reference:
// https://github.com/indeedeng/lsmtree/blob/master/recordlog/src/main/java/com/indeed/lsmtree/recordlog/BasicRecordFile.java
//
package wal
