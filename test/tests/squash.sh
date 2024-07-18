# Test `jit squash` command
set -ex

setup_test_repo squash/repo1
jit clone squash/repo1 --fork=false
cd $JIT_DIR/squash/repo1

# Create new branch
jit new branch1 main
cd branch1
# Add commits to branch
add_commit '1.txt'
add_commit '2.txt'
add_commit '3.txt'
add_commit '4.txt'

# Set up mock response from gh
export GH_RESPONSE='/home/ubuntu/gh-response'
echo '{"baseRefName": "main"}' > $GH_RESPONSE

# Squash commits
jit squash
# Check squashed commit contains all files
git diff --name-only HEAD HEAD~1 | grep '1.txt'
git diff --name-only HEAD HEAD~1 | grep '2.txt'
git diff --name-only HEAD HEAD~1 | grep '3.txt'
git diff --name-only HEAD HEAD~1 | grep '4.txt'
# Check there's only one commit on top of main
[[ $(git log main..HEAD --oneline | wc -l) == 1 ]];

# Cleanup
rm -rf $GIT_PROJECT_ROOT/squash
rm -rf $JIT_DIR/squash
rm -rf $GH_RESPONSE
