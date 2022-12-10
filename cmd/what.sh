#!/bin/bash
# Ask GitHub what this branch is about
# usage:   jit what <branch-name>

# TODO: when given no arguments, show state of all branches
BRANCH=$1
echo "Branch: $BRANCH"

if [[ -z "${GH_USER}" ]]
then
  echo "Unknown GitHub username. Please set the GH_USER environment variable to your GitHub username, then try again."
  exit 1
else
	# TODO: use upstream instead of GH_USER
	# TODO: for some reason 'gh pr view' isn't finding anything?
  gh pr view "$GH_USER:$BRANCH" --json title,state,headRefName,baseRefName,url -t \
  '{{.title}}
{{.state}}: {{.headRefName}} -> {{.baseRefName}}
{{.url}}
'
fi
