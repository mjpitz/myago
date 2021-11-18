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

// Status reports an optional code and message along with the request.
type Status struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Frame is the generalized structure passed along the wire.
type Frame struct {
	Nonce  string      `json:"nonce,omitempty"`
	Status *Status     `json:"status,omitempty"`
	Body   interface{} `json:"body"`
}

type Invoke struct {
	Method string `json:"method,omitempty"`
}
