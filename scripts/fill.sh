#!/usr/bin/env bash

REALPATH=`realpath $0`
DIRPATH=`dirname $REALPATH`

for file in $(ls -1 ../data/excel/); do
./scripts/test.sh parse vehicle xlsx --path ../data/excel/$file
done

for file in $(ls -1 ../data/fs/); do
./scripts/test.sh parse vehicle formstruct --path ../data/fs/$file
done

for file in $(ls -1 ../data/dump/2/); do
./scripts/test.sh parse vehicle fsdump --path ../data/dump/2/$file
done

for file in $(ls -1 ../data/godoc/); do
./scripts/test.sh parse vehicle gsheet --path ../data/godoc/$file
done

for file in $(ls -1 ../data/issued/); do
./scripts/test.sh parse vehicle issued --path ../data/issued/$file
done

for file in $(ls -1 ../data/people/); do
./scripts/test.sh parse people gsheet --path ../data/people/$file
done


DELETE FROM passes;
DELETE FROM bids;
DELETE FROM companies;
DELETE FROM issued;

DELETE FROM passes_people;
DELETE FROM bids_people;
DELETE FROM companies_people;
DELETE FROM issued_people;

ALTER TABLE passes AUTO_INCREMENT = 1;
ALTER TABLE bids AUTO_INCREMENT = 1;
ALTER TABLE companies AUTO_INCREMENT = 1;
ALTER TABLE issued AUTO_INCREMENT = 1;

ALTER TABLE passes_people AUTO_INCREMENT = 1;
ALTER TABLE bids_people AUTO_INCREMENT = 1;
ALTER TABLE companies_people AUTO_INCREMENT = 1;
ALTER TABLE issued_people AUTO_INCREMENT = 1;