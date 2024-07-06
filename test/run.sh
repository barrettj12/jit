#!/bin/bash
# This script runs the E2E tests.
# It should be invoked from the root directory as test/run.sh.
set -e

go build -buildvcs=false -o test/_build/jit .
go build -buildvcs=false -o test/_build/gitserver ./test/gitserver

IMAGE_NAME='jit-test'
CONTAINER_NAME='jit-test'
docker build test --tag $IMAGE_NAME
docker run -dit --name $CONTAINER_NAME $IMAGE_NAME
# shellcheck disable=SC2064
trap "docker rm -f $CONTAINER_NAME >/dev/null" EXIT HUP INT TERM

set +e
for FILE in test/tests/*; do
  [ -e "$FILE" ] || continue
  echo "====== Running test '$(basename $FILE)' ... ==================="

  # We pipe the test to `docker exec` to run inside the container.
  # First we pipe all the 'includes' files (utility functions), followed by the
  # actual test case.
  cat test/includes/* $FILE | docker exec -i $CONTAINER_NAME bash
  RETVAL=$?

  echo
  if [ $RETVAL -eq 0 ]; then
    echo "====== Test '$(basename $FILE)' PASSED ========================"
  else
    echo "====== Test '$(basename $FILE)' FAILED ========================"
    exit 1
  fi
done