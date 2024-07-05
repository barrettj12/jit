#!/bin/bash
# This script runs the E2E tests.
# It should be invoked from the root directory as test/run.sh.
set -ex

go build -buildvcs=false -o test/_build/jit .
go build -buildvcs=false -o test/_build/gitserver ./test/gitserver

docker build test --tag 'jit-test'
docker run -dit --name 'jit-test' 'jit-test'
trap 'docker rm -f jit-test' EXIT

docker exec -it 'jit-test' bash
