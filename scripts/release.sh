#!/usr/bin/env bash

ARG1=$1

REALPATH=`realpath $0`
DIRPATH=`dirname $REALPATH`

cd $DIRPATH/..

set -ex

BRANCH=`git branch | sed -n '/\* /s///p'`

REGISTRY=dockereg.athletic.cloud
IMAGE=itube-meta-$BRANCH

docker run --rm -v "$PWD":/app treeder/bump patch

VERSION=`cat VERSION`
echo "version: $VERSION"

# tag it
git add -A
git commit -m "version $VERSION"
git tag -a "$VERSION" -m "version $VERSION"
git push origin $BRANCH
git push --tags
