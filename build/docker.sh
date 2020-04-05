#!/usr/bin/env bash

REGISTRY_URL=dockereg.athletic.cloud
STAGE=${STAGE:-dev}

VERSION=$(git describe)
IMAGE=${REGISTRY_URL}/collect-proxy-${STAGE}

docker build -t ${IMAGE} .
docker tag ${IMAGE}:latest ${IMAGE}:${VERSION}

docker push ${IMAGE}:${VERSION}
docker push ${IMAGE}:latest
