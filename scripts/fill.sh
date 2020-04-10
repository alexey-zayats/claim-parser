#!/usr/bin/env bash

REALPATH=`realpath $0`
DIRPATH=`dirname $REALPATH`

./scripts/test.sh parse xlsx --path ../data/excel/1.xlsx
./scripts/test.sh parse formstruct --path ../data/fs/1.txt

./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e851500497393b7298b4577
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e85150d497393b7298b4578
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e8515194973933f2a8b4582
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e85151e497393702d8b457b
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e851523497393b7298b4579
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e8515294973933f2a8b4583
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e851530497393702d8b457c
./scripts/test.sh parse fsdump --path ../data/dump/2/form-5e85a1b8497393702d8b458b

./scripts/test.sh parse godoc --path ../data/godoc/1.xlsx
./scripts/test.sh parse godoc --path ../data/godoc/2.xlsx

./scripts/test.sh parse issued --path ../data/issued/1.xlsx
./scripts/test.sh parse issued --path ../data/issued/2.xlsx