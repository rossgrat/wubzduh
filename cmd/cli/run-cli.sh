#!/bin/bash
set -a
. ../../env.txt
set +a
go build .
./cli
