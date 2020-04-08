#!/usr/bin/env bash

REALPATH=`realpath $0`
DIRPATH=`dirname $REALPATH`

binary=claim-parser
env=claim-parser

set -a

if [[ $1 = "-b" ]];
then
    shift 1 ;
    make binary
fi

source $DIRPATH/${env}.env
BINARY=$DIRPATH/../bin/${binary}

$BINARY $@

exit


./scripts/test.sh parse xlsx --path ../data/excel/1.xlsx
./scripts/test.sh parse formstruct --path ../data/fs/1.txt
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e8515194973933f2a8b4582
./scripts/test.sh parse godoc --path ../data/godoc/1.xlsx
./scripts/test.sh parse godoc --path ../data/godoc/2.xlsx

./scripts/test.sh parse registry --path ../data/registry/1.xlsx