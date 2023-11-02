#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

function usage {
    echo "Usage: changelog.bash last-release this-release"
}
trap usage EXIT

lastRelease=$1
thisRelease=$2

lastReleaseSHA=$(git rev-parse $lastRelease)
thisReleaseSHA=$(git rev-parse $thisRelease)

go run github.com/hashicorp/go-changelog/cmd/changelog-build@latest -entries-dir .changelog -last-release $(git rev-parse v0.0.2) -this-release $(git rev-parse HEAD) -changelog-template .changelog/changelog.tmpl -note-template .changelog/note.tmpl
