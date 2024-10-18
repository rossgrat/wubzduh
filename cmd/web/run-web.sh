#!/bin/bash
go build .
set -a
. ../../env.txt
set +a
go build .
./web
