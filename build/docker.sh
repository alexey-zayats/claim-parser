#!/usr/bin/env bash

REALPATH=`realpath $0`
DIRPATH=`dirname $REALPATH`

cd $DIRPATH/..

REGISTRY_URL=dockereg.athletic.cloud
VERSION=$(cat VERSION)
IMAGE=${REGISTRY_URL}/claim-parser

docker build -t ${IMAGE} .
docker tag ${IMAGE}:latest ${IMAGE}:${VERSION}

docker push ${IMAGE}:${VERSION}
docker push ${IMAGE}:latest