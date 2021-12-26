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

package main

// Generate ca.key and ca.crt
//go:generate openssl genrsa -out ca.key 4096
//go:generate openssl req -new -x509 -key ca.key -days 1 -out ca.crt -config ca.conf

// Generate tls.key and tls.csr
//go:generate openssl genrsa -out tls.key 4096
//go:generate openssl req -new -key tls.key -out tls.csr -config tls.conf

// Sign tls.csr using the ca and output to tls.crt
//go:generate openssl x509 -req -in tls.csr -CA ca.crt -CAkey ca.key -CAcreateserial -days 1 -out tls.crt -extensions req_ext -extfile tls.conf

func main() {}
