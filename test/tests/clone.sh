# Test jit clone
set -ex

# Set up test repo to be cloned
setup_test_repo clone/repo1

jit clone clone/repo1 --fork=false

# Test clone is set up correctly
cd $JIT_DIR/clone/repo1
# Test remote is named after user
git remote | grep 'clone'
# Check remote is correctly configured
[ "$(git config remote.clone.fetch)" == '+refs/heads/*:refs/remotes/clone/*' ]
# Test initial worktree 'main' was created
ls | grep 'main'
[ "$(git rev-parse --abbrev-ref 'main@{u}')" == 'clone/main' ]

# Cleanup
rm -rf $GIT_PROJECT_ROOT/clone
rm -rf $JIT_DIR/clone
