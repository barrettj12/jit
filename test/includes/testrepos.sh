# Sets up a basic test repo in the specified location
#   usage:  setup_test_repo <user>/<repo>
setup_test_repo() {
  if [[ $GIT_PROJECT_ROOT == "" ]]; then
    echo "env var GIT_PROJECT_ROOT not defined"
    exit 1
  fi

  TEST_REPO_PATH="$GIT_PROJECT_ROOT/$1"
  mkdir -p "$TEST_REPO_PATH"
  cd $TEST_REPO_PATH

  # Initialise Git repo
  git init
  # Create a file
  echo "hello world" > foo.txt
  # git add and commit
  git add foo.txt
  git commit -m "Initial commit"
}
