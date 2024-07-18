# Test `jit rebase` command
set -ex

setup_test_repo rebase/repo1
jit clone rebase/repo1 --fork=false
cd $JIT_DIR/rebase/repo1

# Create new branch
jit new branch1 main
cd branch1
# Add commits to branch
add_commit 'new.txt'
# Add commits to remote branch
add_remote_commit rebase/repo1 'base1.txt'
add_remote_commit rebase/repo1 'base2.txt'
add_remote_commit rebase/repo1 'base3.txt'

# Set up mock response from gh
export GH_RESPONSE='/home/ubuntu/gh-response'
echo '{"baseRefName": "main"}' > $GH_RESPONSE

# Rebase
jit rebase
# Check commits are in the correct order in the git log
git show -s --oneline HEAD | grep 'new.txt'
git show -s --oneline HEAD~1 | grep 'base3.txt'
git show -s --oneline HEAD~2 | grep 'base2.txt'
git show -s --oneline HEAD~3 | grep 'base1.txt'

# Cleanup
rm -rf $GIT_PROJECT_ROOT/rebase
rm -rf $JIT_DIR/rebase
rm -rf $GH_RESPONSE
