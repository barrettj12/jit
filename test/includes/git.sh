# Adds a file in the current directory and commits it
#   usage:  add_commit <filename>
add_commit() {
  FILENAME=$1
  echo "goodbye cruel world" > $FILENAME
  git add $FILENAME
  git commit -m "add $FILENAME"
}

# Adds a file to the given remote and commits it
#   usage:  add_commit <remote> <filename>
add_remote_commit() {(
  cd "$GIT_PROJECT_ROOT/$1"
  add_commit $2
)}

# Adds a branch to the given remote and commits it
#   usage:  add_commit <remote> <new-branch> <base>
add_remote_branch() {(
  cd "$GIT_PROJECT_ROOT/$1"
  git branch $2 $3
)}
