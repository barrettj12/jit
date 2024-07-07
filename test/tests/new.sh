# Test `jit new` command
set -ex

# Set up test repo to be cloned
setup_test_repo new/repo1
jit clone new/repo1 --fork=false
cd $JIT_DIR/new/repo1

# Test creating new branch based on existing branch
git branch branch1 main
jit new branch1
cd branch1
# shellcheck disable=SC2063
git branch | grep '* branch1'
cd ..

# Test creating new worktree from existing local branch
jit new branch2 main
cd branch2
# shellcheck disable=SC2063
git branch | grep '* branch2'
cd ..

# Test creating new worktree from existing remote branch
# Clone the test repo
git clone "$GIT_PROJECT_ROOT/new/repo1" "$GIT_PROJECT_ROOT/new2/repo1"
git -C "$GIT_PROJECT_ROOT/new2/repo1" branch branch3 main
# Don't add the new remote - the new command should automatically do this
jit new new2:branch3
cd branch3
# shellcheck disable=SC2063
git branch | grep '* branch3'
# Check upstream is set
git rev-parse --abbrev-ref 'branch3@{u}' | grep 'new2/branch3'
git rev-parse --abbrev-ref 'branch3@{push}' | grep 'new2/branch3'
cd ..
