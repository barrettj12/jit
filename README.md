# Jit - *the superior way to use Git*

Jit is an alternative "porcelain" (wrapper/interface) for Git. It is
intentionally opinionated, and aims to push the user to use Git in a very
particular way, which the developer, at least, believes to be a superior way.

Some of the ways Jit differs from vanilla Git include:
- Worktrees as the standard workflow
- Baked-in support for triangular workflows (new branches are pushed to the
  user's fork by default)
- There are no special remote names such as "origin" or "upstream" - all
  remotes are named after the corresponding GitHub user.

It also includes some convenience commands, such as being able to apply a patch
directly from GitHub using `jit apply`.

At the same time, Jit works within the language of Git - its repos are Git
repos, its branches Git branches, and so on. This means you can still run
standard Git commands in any repo or branch set up using Jit.

**NB: although I use Jit extensively, it is still in development, and may
naturally contain bugs and incomplete features. Please report issues
[here](https://github.com/barrettj12/jit/issues), or even better,
[open a PR](https://github.com/barrettj12/jit/pulls).**

## Dependencies
- [Git](https://git-scm.com/)
- [GitHub CLI](https://cli.github.com/) (`gh`, for certain commands)

## Installation
```
go install github.com/barrettj12/jit@latest
```

## The Jit workflow

### Directory structure

Jit uses the `JIT_DIR` environment variable to determine where/how to store
Git repos. Inside this directory, it uses the following file structure:
```
${JIT_DIR}
├─ user1
│  ├─ repo1
│  │  ├─ branch1
│  │  ├─ branch2
│  │  ├─ ...
│  ├─ repo2
│  │  ├─ ...
├─ user2
│  ├─ repo3
│  ├─ ...
├─ ...
```

For example, say we have a repo located at
```
https://github.com/johnsmith/foo-app
```
with branches `main`, `stable` and `exp`. When you run
```
jit clone https://github.com/johnsmith/foo-app
```
Jit will create the following file structure:
```
${JIT_DIR}
└─ johnsmith
   └─ foo-app
      ├─ main
      ├─ stable
      └─ exp
```

Each subdirectory of `foo-app` is a worktree (see below) with the corresponding
branch checked out.

### Worktrees

Worktrees are a great feature of Git added in later versions. They allow you to
check out multiple branches at once. Jit uses worktrees as the default workflow.
This makes it easy to work on multiple things simultaneously - just move to a
different directory, instead of fiddling around with `git stash`, `git switch`,
etc.

### Triangular workflows

Modern software development (in teams) often uses a "triangular workflow" with
three copies of the source code:
- A central "source" copy
- A remote fork
- A local copy

Developers pull changes from the source copy to their local copy, then push
changes to their fork. After this, a pull request is opened to merge changes
from their fork into the central source.

Git does support this kind of workflow, but it can be a lot of work to set it
up and maintain it. Jit makes it easy by automatically setting up the
remotes - Jit commands are designed to use a triangular workflow.

All remotes are named after the GitHub user by default. So, when you clone
`https://github.com/ecma/gizmo`, the "origin" remote will be named `ecma`,
while your personal fork will be the same as your GitHub username.

## Jit commands

### `jit clone`

Clone a GitHub repository. Jit sets up the local repository with a bare/worktree
structure, so that other Jit commands can use it. 
```
jit clone <user> <repo>
jit clone <user>/<repo>
jit clone https://github.com/<user>/<repo>
```

### `jit new`
Create a new worktree/branch. There are three different modes:
- `jit new <branch>` - create a new worktree for an existing branch
- `jit new <branch> <base>` - create a new branch based on an existing branch
- `jit new <remote>:<branch>` - check out a branch from someone else's remote

### `jit pull`
Update a given branch.
```
jit pull <branch>
```

### `jit rm`
Remove a given branch/worktree and its upstream remote.
```
jit rm <branch>
```

### `jit squash`
Squash all commits in the current branch into one. Requires a GitHub PR to be
open on the branch.

### `jit rebase`
Rebase the current branch against the latest version of the base branch.
Requires a GitHub PR to be open on the branch.

<!--
TODO: restore this section once hooks are implemented

## Hooks
*This has not yet been implemented.*

Hooks are Bash scripts which tell Jit to perform extra actions on some command.
For example, say when I create a new branch:
```
jit branch new old
```
I want some build artifacts to be copied from `old` to `new`. In this case,
I can write a Bash script and put it in the repo folder (which contains the
branches):

`branch.sh`:
```bash
#!/bin/bash
cp $OLD/_build $NEW/_build
```
-->

## History

Jit started from a series of Bash scripts I wrote to automate my Git workflow.
At some point, I realised these scripts could be turned into a fully-fledged
CLI tool.
