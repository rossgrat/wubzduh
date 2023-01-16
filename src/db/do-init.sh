#!/bin/bash
set -a
. ../../src/env.txt
set +a
psql -U $DB_USERNAME -d $DB_NAME -a -f init-db.sql
