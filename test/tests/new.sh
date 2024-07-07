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
# First add a commit to the remote's 'main' branch - we want to check the local
# 'main' branch is updated before the new branch is created.
echo "goodbye cruel world" > "$GIT_PROJECT_ROOT/new/repo1/bar.txt"
git -C "$GIT_PROJECT_ROOT/new/repo1" add bar.txt
git -C "$GIT_PROJECT_ROOT/new/repo1" commit -m "remote file added"
# Create the new branch
jit new branch2 main
cd branch2
# shellcheck disable=SC2063
git branch | grep '* branch2'
git log | grep 'remote file added'
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


# Check that branch names containing a slash '/' are supported
# Regression test for https://github.com/barrettj12/jit/issues/6
git branch 'branch/with/slashes' main
jit new 'branch/with/slashes'
cd 'branch_with_slashes'
# shellcheck disable=SC2063
git branch | grep '* branch/with/slashes'
cd ..
