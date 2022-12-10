# Jit

Jit is an alternative "porcelain" (wrapper/interface) for Git.

**NB: although I use Jit extensively, it is very much a work-in-progress
and not yet suitable for general release. Contributions are appreciated.**

## Dependencies
- [Git](https://git-scm.com/)
- [GitHub CLI](https://cli.github.com/) (`gh`)

## Installation
1. Install the above dependencies.
2. Download the source code.
3. Add the source tree root directory (containing the downloaded code) to 
   your PATH.

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
│  ├─ repo1
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

We opt for remotes called `source` and `fork` instead, as we find the Git terms
`upstream` and `origin` to be unclear.

*TODO: thinking of naming all remotes based on the GitHub user instead.*

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
Create a new worktree/branch.
```
jit new <branch> [<remote>/]<base>
```

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

### `jit what`
Ask Jit what a given branch is about.
```
jit what <branch>
```

### `jit where`
Ask Jit where the given branch will push to / pull from.
```
jit where push <branch>
jit where pull <branch>
```

Any commands unrecognised by Jit will be sent verbatim to Git. This allows
`jit` to be used as a drop-in replacement for `git`.

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

## History

Jit started from a series of Bash scripts I wrote to automate my Git workflow.
At some point, I realised these scripts could be turned into a fully-fledged
CLI tool.