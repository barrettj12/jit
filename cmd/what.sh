#!/bin/bash
# Ask GitHub what this branch is about
# usage:   jit what <branch-name>

BRANCH=$1

if [[ -z "${GH_USER}" ]]
then
  echo "Unknown GitHub username. Please set the GH_USER environment variable to your GitHub username, then try again."
  exit 1
else
  gh pr view "$GH_USER:$BRANCH" --json title,state,headRefName,baseRefName,url -t \
  '{{.title}}
{{.state}}: {{.headRefName}} -> {{.baseRefName}}
{{.url}}
'
fi
