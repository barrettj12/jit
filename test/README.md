# Testing
This folder contains tests for Jit. The tests are "end-to-end" style tests -
they use the compiled `jit` binary along with a real version of Git. The only
thing that is mocked out is GitHub, using a fake Git server (implementation can
be seen in the `gitserver` folder).

To run the test suite, invoke `test/run.sh` from the root directory. This is a
Bash script that spins up a Docker container and runs the tests inside this
container. We use `git config` to rewrite the github.com URLs into localhost
ones, which redirect to a fake Git server running inside the container.

Test cases are stored as Bash scripts in the `tests` subfolder. Each test case
is run by piping the script into `docker exec` (along with all the testing
utilities defined in the `includes` folder).

The structure of this folder:
- `_build`: this is where the built binaries are placed (`jit` and `gitserver`)
    before being copied into the Docker container.
- `gitserver`: contains a fake Git server implementation which is used to mock
    out GitHub in the tests.
- `includes`: contains Bash files which define utility functions to use in
    tests.
- `tests`: contains test cases as Bash files. These are piped into
    `docker exec` and run inside the container.
- `Dockerfile`: defines the Docker container which tests are run inside.
- `run.sh`: the entrypoint to the test suite.
