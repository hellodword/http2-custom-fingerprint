#!/usr/bin/env bash

set -e
set -x

tag="$1"
[ -n "$tag" ]

revision_file="$2"
[ -f "$revision_file" ]

branch="$tag-patched"

if grep "$tag" "$revision_file"; then
    echo "$tag already done"
    exit 0
fi

sha="$(git rev-parse --verify "origin/$branch")"
[ -n "$sha" ]

mkdir -p /tmp/net-revision
pushd /tmp/net-revision
go mod init x
revision="$((go get -u "github.com/hellodword/http2-custom-fingerprint@$sha" 2>&1 || true)  | grep -oP 'v\d+\.\d+\.\d+-\d+\.\d+\-[\da-f]{12}' | head -1)"
popd
rm -rf /tmp/net-revision

[ -n "$revision" ]

echo "replace golang.org/x/net $tag => github.com/hellodword/http2-custom-fingerprint $revision" | tee -a "$revision_file"
