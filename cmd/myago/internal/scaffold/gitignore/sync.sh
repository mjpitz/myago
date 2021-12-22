#!/usr/bin/env sh

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
