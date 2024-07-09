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

# Cleanup
rm -rf $GIT_PROJECT_ROOT/pull
rm -rf $JIT_DIR/pull
