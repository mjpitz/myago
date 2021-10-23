#!/usr/bin/env bash

for pkg in $(go list ./...); do
  output="${pkg#"github.com/mjpitz/myago/"}/README.md"

  godocdown "$pkg" > "$output"
done
