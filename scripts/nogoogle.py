#!/usr/bin/env python3

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
            parts = line.split(" ")
            module = parts[0]

            cwd = os.getcwd()
            if key:
                cwd = cwd + '/' + key

            result = subprocess.run(["go", "mod", "why", module], cwd=cwd, stdout=subprocess.PIPE, universal_newlines=True)

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
                documented = module in allowed
                if not documented:
                    undocumented.append({
                        "key": key,
                        "module": module,
                    })

                discovered.append({
                    "key": key,
                    "module": module,
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
        f"key:    '{dependency['key']}'\n"
        f"module: '{dependency['module']}'\n"
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
            f"\tkey:    '{dependency['key']}'"
            f"\tmodule: '{dependency['module']}'"
        )
    exit(1)
