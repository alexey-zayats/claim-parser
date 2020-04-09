#!/usr/bin/env bash

REALPATH=`realpath $0`
DIRPATH=`dirname $REALPATH`

FILES="
drop.sql
districts.sql
users.sql
files.sql
bids.sql
passes.sql
issued.sql
"

for file in $FILES; do
    echo $DIRPATH/$file
    mysql -u pass pass < $DIRPATH/$file
done