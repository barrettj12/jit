# Test jit clone
set -ex

# Set up test repo to be cloned
setup_test_repo clone/repo1

jit clone clone/repo1 <<< 'n' # don't create a fork

# Test clone is set up correctly
cd $JIT_DIR/clone/repo1
# Test remote is named after user
git remote | grep 'clone'
# Test initial worktree 'main' was created
ls | grep 'main'
