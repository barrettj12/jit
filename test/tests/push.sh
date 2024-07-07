# Test `jit push` command
set -ex

setup_test_repo push/repo1
jit clone push/repo1 --fork=false
cd $JIT_DIR/push/repo1
export GH_USER='push'


# Create branch - initially should not have push target set up
jit new branch1 main
cd branch1
# Check no upstream is set
[ -z "$(git rev-parse --abbrev-ref '@{push}')" ]
# Add a commit and push it
add_commit one.txt
jit push
 # Check commit exists on remote
git -C "$GIT_PROJECT_ROOT/push/repo1" log branch1 | grep 'one.txt'
git rev-parse --abbrev-ref '@{push}' || grep 'push/branch1'
