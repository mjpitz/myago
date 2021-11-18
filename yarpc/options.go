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

package yarpc

import (
	"context"
	"crypto/tls"

	"github.com/hashicorp/yamux"
	"github.com/jonboulle/clockwork"

	"github.com/mjpitz/myago/encoding"
)

// Option defines an generic way to configure clients and servers.
type Option func(opt *options)

type options struct {
	context  context.Context
	yamux    *yamux.Config
	tls      *tls.Config
	encoding *encoding.Encoding
	clock    clockwork.Clock
}

// WithTLS enables TLS.
func WithTLS(config *tls.Config) Option {
	return func(opt *options) {
		if config != nil {
			opt.tls = config
		}
	}
}

// WithYamux configures yamux using the provided configuration.
func WithYamux(config *yamux.Config) Option {
	return func(opt *options) {
		if config != nil {
			opt.yamux = config
		}
	}
}

// WithContext provides a custom context to the underlying system. Mostly used on servers.
func WithContext(ctx context.Context) Option {
	return func(opt *options) {
		if ctx != nil {
			opt.context = ctx
		}
	}
}

// WithEncoding configures how messages are serialized.
func WithEncoding(encoding *encoding.Encoding) Option {
	return func(opt *options) {
		if encoding != nil {
			opt.encoding = encoding
		}
	}
}

// private options used to configure the underlying stream... mainly to keep the wrap() interface reasonable

type streamOption func(stream *rpcStream)

func withContext(context context.Context) streamOption {
	return func(stream *rpcStream) {
		if context != nil {
			stream.context = context
		}
	}
}

func withEncoding(encoding *encoding.Encoding) streamOption {
	return func(stream *rpcStream) {
		if encoding != nil {
			stream.encoding = encoding
		}
	}
}
