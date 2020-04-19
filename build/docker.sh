#!/usr/bin/env bash

REALPATH=`realpath $0`
DIRPATH=`dirname $REALPATH`

cd $DIRPATH/..

REGISTRY_URL=dockereg.athletic.cloud
REGISTRY_URL=aazayats

IMAGE=${REGISTRY_URL}/claim-parser
VERSION=$(cat VERSION)

docker build -t ${IMAGE} .
docker tag ${IMAGE}:latest ${IMAGE}:${VERSION}

docker push ${IMAGE}:${VERSION}
docker push ${IMAGE}:latest