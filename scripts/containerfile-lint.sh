#!/bin/bash

set -eu

# Fetch hadolint
version="2.12.0"
target="hadolint"
url="https://github.com/${target}/${target}/releases/download/v${version}/${target}-Linux-x86_64"

WORK_DIR=$(mktemp -d)
target_path="${WORK_DIR}/${target}"
trap 'rm -rfv ${WORK_DIR} &>/dev/null' EXIT

curl -L "${url}" -o "${target_path}"
chmod +x "${target_path}"

FILES=$(find . -name 'Containerfile*' -o -name 'Dockerfile*')

while read -r file; do
  echo "Linting: ${file}"
  # Doesn't support specifying multiple files, see https://github.com/hadolint/hadolint-action/issues/3
  ${target_path} "${file}"
done <<< "${FILES}"
