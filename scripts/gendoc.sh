#!/usr/bin/env bash

for pkg in $(go list ./... | egrep -v '^github.com/mjpitz/myago$'); do
  output="${pkg#"github.com/mjpitz/myago/"}/README.md"

  godocdown "$pkg" > "$output"
done
