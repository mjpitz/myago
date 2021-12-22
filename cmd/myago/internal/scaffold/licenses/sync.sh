#!/usr/bin/env sh

license_sync() {
  name="${1}"

  wget -qO "${name}.txt" https://raw.githubusercontent.com/licenses/license-templates/master/templates/${name}.txt
  wget -qO "${name}-header.txt" https://raw.githubusercontent.com/licenses/license-templates/master/templates/${name}-header.txt || {
    rm "${name}-header.txt"
  }
}

license_sync "agpl3"
license_sync "mit"
license_sync "mpl"
