#!/usr/bin/env bash

set -e
set -x

# git rev-parse --verify 

tag="$1"
[ -n "$tag" ]
patch_file="$2"
[ -f "$patch_file" ]

branch="$tag-patched"

if git rev-parse --verify "origin/$branch"; then
  echo "$tag already patched"
  exit 0
fi


git checkout "$tag"
git checkout -b "$branch"
git apply "$patch_file"
git add .
git commit -m 'patch http2 fingerprinting'
