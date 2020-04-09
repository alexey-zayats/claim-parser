#!/usr/bin/env bash

./scripts/test.sh parse xlsx --path ../data/excel/1.xlsx
./scripts/test.sh parse formstruct --path ../data/fs/1.txt
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e8515194973933f2a8b4582
./scripts/test.sh parse godoc --path ../data/godoc/1.xlsx
./scripts/test.sh parse godoc --path ../data/godoc/2.xlsx

./scripts/test.sh parse issued --path ../data/issued/1.xlsx
./scripts/test.sh parse issued --path ../data/issued/2.xlsx