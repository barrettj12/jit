# Test `jit edit` command
set -ex

# Set up test repo
setup_test_repo edit/repo1
jit clone edit/repo1 --fork=false
cd $JIT_DIR/edit/repo1

# Create a new branch and open it for editing
jit new branch1 main --no-edit
jit edit branch1

# Editing a branch that doesn't exist should cause an error
OUTPUT=$( ! jit edit nonexistent 2>&1 )
echo $OUTPUT | grep 'no worktree found for branch "nonexistent"'

# edit accepts 'remote:branch'
jit edit user:branch1

# Cleanup
rm -rf $GIT_PROJECT_ROOT/new
rm -rf $GIT_PROJECT_ROOT/new2
rm -rf $JIT_DIR/new
