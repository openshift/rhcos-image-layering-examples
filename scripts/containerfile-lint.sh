#!/bin/bash

set -eu

FILES=$(find . -name 'Containerfile*' -o -name 'Dockerfile*')

while read -r file; do
  echo "Linting: ${file}"
  # Doesn't support specifying multiple files, see https://github.com/hadolint/hadolint-action/issues/3
  podman run --rm -i ghcr.io/hadolint/hadolint < "${file}"
done <<< "${FILES}"
