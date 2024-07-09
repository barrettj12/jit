# Test `jit where` command
set -ex

setup_test_repo where/repo1
jit clone where/repo1 --fork=false
cd $JIT_DIR/where/repo1
export GH_USER='where2'


# Test push/pull targets for main
jit where push main | grep 'where/main'
jit where pull main | grep 'where/main'

# Test `where` inside worktree with no arg
cd main
jit where push | grep 'where/main'
jit where pull | grep 'where/main'
cd ..


# Create a fork for the user
git clone "$GIT_PROJECT_ROOT/where/repo1" "$GIT_PROJECT_ROOT/where2/repo1"

# Test push/pull targets are not set for new branch
jit new branch1 main
jit where push branch1 | grep 'no push target configured'
jit where pull branch1 | grep 'no pull target configured'

# Test that pushing sets the targets for the branch
cd branch1
jit push
jit where push | grep 'where2/branch1'
jit where pull | grep 'where2/branch1'
cd ..


# Cleanup
rm -rf $GIT_PROJECT_ROOT/where
rm -rf $GIT_PROJECT_ROOT/where2
rm -rf $JIT_DIR/where
