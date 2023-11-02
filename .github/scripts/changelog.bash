#!/usr/bin/env bash

lastRelease=$1
thisRelease=$2

lastReleaseSHA=$(git rev-parse $1)
thisReleaseSHA=$(git rev-parse $2)

go run github.com/hashicorp/go-changelog/cmd/changelog-build@latest -entries-dir .changelog -last-release $(git rev-parse v0.0.2) -this-release $(git rev-parse HEAD) -changelog-template .changelog/changelog.tmpl -note-template .changelog/note.tmpl
