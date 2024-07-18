#!/bin/bash
# This script runs the E2E tests.
# It should be invoked from the root directory as test/run.sh.
#   test/run.sh        runs all tests
#   test/run.sh foo    runs only test cases matching 'foo'
set -e

# Optional argument is a string defining which test cases to run.
RUN_TEST_SPEC=$1

go build -buildvcs=false -o test/_build/jit .
go build -buildvcs=false -o test/_build/gitserver ./test/gitserver
go build -buildvcs=false -o test/_build/goland ./test/goland
go build -buildvcs=false -o test/_build/gh ./test/gh

IMAGE_NAME='jit-test'
CONTAINER_NAME='jit-test'
docker build test --tag $IMAGE_NAME
docker run -dit --name $CONTAINER_NAME $IMAGE_NAME
# shellcheck disable=SC2064
trap "docker rm -f $CONTAINER_NAME >/dev/null" EXIT HUP INT TERM

set +e
for FILE in test/tests/*; do
  [ -e "$FILE" ] || continue
  CURRENT_TEST=$(basename $FILE)

  # If a specific test was requested, check if this matches
  if [ -n "$RUN_TEST_SPEC" ] && ! echo $CURRENT_TEST | grep $RUN_TEST_SPEC; then
    continue
  fi

  echo "====== Running test '$CURRENT_TEST' ... ==================="

  # We pipe the test to `docker exec` to run inside the container.
  # First we pipe all the 'includes' files (utility functions), followed by the
  # actual test case.
  cat test/includes/* $FILE | docker exec -i $CONTAINER_NAME bash
  RETVAL=$?

  echo
  if [ $RETVAL -eq 0 ]; then
    echo "====== Test '$CURRENT_TEST' PASSED ========================"
  else
    echo "====== Test '$CURRENT_TEST' FAILED ========================"
    exit 1
  fi
done
