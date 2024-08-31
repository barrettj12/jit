# Test `jit squash` command
set -ex

setup_test_repo squash/repo1
jit clone squash/repo1 --fork=false
cd $JIT_DIR/squash/repo1

# This is a hack to stop git complaining about the lack of editor in this env.
# It will fill in a default squashed commit message for us.
export EDITOR='true'

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
cd ..

# Regression test for https://github.com/barrettj12/jit/issues/57
jit new branch2 branch1
cd main
add_remote_commit squash/repo1 '1.txt' 'conflict'
cd ../branch2
# pull new commit from main, merge it into branch2 and fix the conflicts
jit pull main
git merge main -s ours
# squash should be able to handle previous conflict resolution
jit squash

# Cleanup
rm -rf $GIT_PROJECT_ROOT/squash
rm -rf $JIT_DIR/squash
rm -rf $GH_RESPONSE
