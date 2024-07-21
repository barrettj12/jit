# Test `jit pull` command
set -ex

# Setup test repo
setup_test_repo pull/repo1
jit clone pull/repo1 --fork=false
cd $JIT_DIR/pull/repo1

# Test we can pull a branch from another by providing the branch name
jit new branch1 main
add_remote_commit pull/repo1 'one.txt'
cd branch1
jit pull main
git log main | grep 'one.txt'
cd ..

# Test we can pull a branch from the repo root
add_remote_commit pull/repo1 'two.txt'
jit pull main
git log main | grep 'two.txt'

# Test we can pull the current branch with no arg
add_remote_commit pull/repo1 'three.txt'
cd main
jit pull
git log main | grep 'three.txt'
cd ..

# Regression test for https://github.com/barrettj12/jit/issues/41
# Create remote branch
add_remote_branch pull/repo1 branch2 main
# Check branch2 is not in cloned repo
git branch -r | grep 'pull/branch2' && exit 1
# Pull main branch
jit pull main
# Check branch2 is still not in cloned repo
git branch -r | grep 'pull/branch2' && exit 1

# Regression test for https://github.com/barrettj12/jit/issues/46
# Create new remote branch
add_remote_branch pull/repo1 branch3 main
# Pull remote branch
jit new pull:branch3
# Delete remote tracking branch
git branch -D --remote pull/branch3
# Try to pull branch
jit pull branch3

# Cleanup
rm -rf $GIT_PROJECT_ROOT/pull
rm -rf $JIT_DIR/pull
