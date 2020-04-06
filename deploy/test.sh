#!/usr/bin/env bash

BIND=/Users/alexis/workspace/quartex/claim-parser/data
ENV=/Users/alexis/workspace/quartex/claim-parser/deploy/env/claim-parser.env

docker run -d --restart always --name claim-parser --env-file $ENV -v $BIND:/data aazayats/claim-parser:latest watch