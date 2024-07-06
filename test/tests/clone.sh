# Test jit clone
set -e

# Set up test repo to be cloned
setup_test_repo clone/repo1

jit clone clone/repo1 <<< 'n' # don't create a fork

# TODO: test clone is set up correctly