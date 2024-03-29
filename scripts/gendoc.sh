#!/usr/bin/env bash
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


for pkg in $(go list all | egrep '^go.pitz.tech/lib' | egrep -v '^go.pitz.tech/lib$' | awk '{print $1}'); do
  output="${pkg#"go.pitz.tech/lib/"}/README.md"
  echo $output

  godocdown -template ./templates/docs/README.md.tmpl "$pkg" > "$output"
done

prettier --write --parser markdown **/README.md
