# Test `jit new` command
set -ex

# Set up test repo to be cloned
setup_test_repo new/repo1
jit clone new/repo1 --fork=false

# Test creating new branch based on existing
# TODO

# Test creating new worktree from existing local branch
# TODO

# Test creating new worktree from existing remote branch
# TODO
