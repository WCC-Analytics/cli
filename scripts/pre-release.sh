#!/bin/bash

set -e

if [ ! -z "$(git status --porcelain)" ]; then
  echo "The git workspace must be clean"
  exit 1
fi

if [ $# -eq 0 ]; then
  echo "You must specify a version"
  exit 1
fi

TAGNAME=$1
if [ "${TAGNAME:0:1}" == "v" ]; then
  echo "Version is incorrect; must match format X.Y.Z"
  exit 1
fi

VERSION="v$TAGNAME"
echo "Using version $VERSION"

# get git in order
git push
git tag -a "$VERSION" -m "$VERSION"
git push origin "$VERSION"  # push only this tag
