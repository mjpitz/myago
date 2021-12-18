#!/usr/bin/env python3
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


# Why? In my opinion (and from some of my experience), Google has shown a blatant disregard with regards to common
# social and ethical issues. So long as they continue to act this way, I am choosing to minimize their technology stack
# in the things that I build.

import glob
import os
import subprocess

allowed_google_libraries = {
    "": {
        "github.com/google/btree": "",
        "google.golang.org/protobuf": "",
    },
    "cmd/myago": {},
    "paxos": {
        "github.com/google/flatbuffers": "",
        "github.com/google/btree": "",
        "google.golang.org/protobuf": "",
    },
}

base = subprocess.run(["go", "list", "-m"], stdout=subprocess.PIPE, universal_newlines=True).stdout.strip()

discovered = []
undocumented = []

for file in glob.iglob("**/go.mod", recursive=True):
    key = file.rstrip("go.mod")
    key = key.rstrip("/")

    allowed = allowed_google_libraries[key]

    handle = open(file, 'r')
    for line in handle:
        line = line.strip()
        if "google" in line:
            dependency = line.split(" ")[0]

            module = base
            cwd = os.getcwd()

            if len(key) > 0:
                module = module + '/' + key
                cwd = cwd + '/' + key

            result = subprocess.run(["go", "mod", "why", dependency], cwd=cwd, stdout=subprocess.PIPE, universal_newlines=True)

            usages = []
            current = []

            for link in result.stdout.split('\n'):
                if 'module does not need package' in link:
                    break
                elif '#' in link:
                    if len(current) > 0:
                        usages.append(" => ".join(current).rstrip(" => "))
                    current = []
                elif len(link) > 0:
                    current.append(link)

            if len(current) > 0:
                usages.append(" => ".join(current).rstrip(" => "))

            if len(usages) > 0:
                documented = dependency in allowed
                if not documented:
                    undocumented.append({
                        "module": module,
                        "dependency": dependency,
                    })

                discovered.append({
                    "module": module,
                    "dependency": dependency,
                    "documented": documented,
                    "usages": usages,
                })

print("# Documented usages:")
for dependency in discovered:
    usages = ""
    for usage in dependency["usages"]:
        usages = usages + "  - " + usage + "\n"

    print(
        f"---\n"
        f"module:     '{dependency['module']}'\n"
        f"dependency: '{dependency['dependency']}'\n"
        f"usages:\n"
        f"{usages}"
    )

if undocumented:
    print("")
    print("")
    print("# The following is a list of undocumented dependencies. Update scripts/nogoogle.py with new dependencies.")
    for dependency in undocumented:
        print(
            f"---"
            f"undocumented:"
            f"\tmodule:    '{dependency['module']}'"
            f"\tdependency: '{dependency['dependency']}'"
        )
    exit(1)
