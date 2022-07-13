#!/bin/bash
# remove a branch (local & remote)

BRANCH=$1
FLAGS=( "${@:2}" )

# Confirm removal
read -rp "Type 'y' to confirm removal of $BRANCH: " RESP
if [ "$RESP" = 'y' ]
then
  # TODO: check before running these commands
  # remove branch in remote fork
  git push fork --delete "$BRANCH"
  # remove worktree
  git worktree remove "$BRANCH" "${FLAGS[@]}"
  # remove branch locally
  git branch -D "$BRANCH"
  # remove folder if it still exists
else
  echo "Abort deletion."
fi
