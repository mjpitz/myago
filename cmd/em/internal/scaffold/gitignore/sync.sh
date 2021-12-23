#!/usr/bin/env sh
# Copyright (C) 2021 Mya Pitzeruse
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published
# by the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.


gitignore_sync() {
  name="${1}"

  wget -qO "${name}.gitignore" https://raw.githubusercontent.com/github/gitignore/main/${name}.gitignore
}

gitignore_sync "Go"
gitignore_sync "Node"
gitignore_sync "Rust"
gitignore_sync "Python"
gitignore_sync "C"
gitignore_sync "C++"
