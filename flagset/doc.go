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

/*
Package flagset provides an opinionated approach to constructing an applications' configuration using Golang structs and
tags. It's designed in a way that allows configuration to be loaded from files, environment variables, and/or command
line flags. The following details the various tags that can be specified on a primitive field.

- `json` - `string` - Configure the name of the flag. Convention is to use snake case.

- `usage` - `string` - Configure the description string of the flag.

- `default` - `any` - Configure the default value for the flag. Can be overridden by setting the value on the struct.

- `hidden` - `bool` - Hides the flag from output. The value can still be configured.

- `required` - `bool` - Specifies that the flag must be specified.

Nested structures are supported, making application configuration composable and portable between systems.

*/
package flagset
