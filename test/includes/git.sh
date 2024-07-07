# Adds a file in the current directory and commits it
#   usage:  add_commit <filename>
add_commit() {
  FILENAME=$1
  echo "goodbye cruel world" > $FILENAME
  git add $FILENAME
  git commit -m "add $FILENAME"
}
