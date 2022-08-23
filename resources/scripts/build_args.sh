#!/usr/bin/env bash

# Skip branch/revision with missing git repo.
if [[ -d .git ]] || git rev-parse --git-dir >/dev/null 2>&1; then
	branch=$(git symbolic-ref HEAD 2>/dev/null)
	[[ -z "$VERSION" ]] && VERSION=$(git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
	revision=$(git log -1 --pretty=format:"%H" 2>/dev/null)
fi

buildUser=${GITHUB_ACTOR:-${USER:-$(whoami)}}
buildDate=$(date +%FT%T%Z)
versionPkg=github.com/hellofresh/hfkcat/internal/version

echo -X "$versionPkg".version="$VERSION" -X "$versionPkg".branch="$branch" -X "$versionPkg".revision="$revision" -X "$versionPkg".buildUser="$buildUser" -X "$versionPkg".buildDate="$buildDate"
