package cmd

import (
	"github.com/barrettj12/jit/common/testutil"
	"os"
	"testing"
)

func TestSquashEditor(t *testing.T) {
	// Set up rebase to-do file
	file, err := os.CreateTemp("", "git-rebase-todo")
	testutil.CheckErr(t, err)
	t.Cleanup(func() {
		err = os.Remove(file.Name())
		testutil.CheckErr(t, err)
	})

	_, err = file.WriteString(squashFileBefore)
	testutil.CheckErr(t, err)
	err = file.Close()
	testutil.CheckErr(t, err)

	err = SquashEditor(nil, []string{file.Name()})
	testutil.CheckErr(t, err)
	fileContents, err := os.ReadFile(file.Name())
	testutil.AssertEqual(t, string(fileContents), squashFileAfter)
}

var squashFileBefore = `
pick bd64edd commit message 1
pick 86304e9 commit message 2
pick dab68cf commit message 3
pick 0c4ebf0 commit message 4

# Rebase a422b17..0c4ebf0 onto a422b17 (4 commands)
#
# Commands:
# p, pick <commit> = use commit
# r, reword <commit> = use commit, but edit the commit message
# e, edit <commit> = use commit, but stop for amending
# s, squash <commit> = use commit, but meld into previous commit
# f, fixup [-C | -c] <commit> = like "squash" but keep only the previous
#                    commit's log message, unless -C is used, in which case
#                    keep only this commit's message; -c is same as -C but
#                    opens the editor
# x, exec <command> = run command (the rest of the line) using shell
# b, break = stop here (continue rebase later with 'git rebase --continue')
# d, drop <commit> = remove commit
# l, label <label> = label current HEAD with a name
# t, reset <label> = reset HEAD to a label
# m, merge [-C <commit> | -c <commit>] <label> [# <oneline>]
#         create a merge commit using the original merge commit's
#         message (or the oneline, if no original merge commit was
#         specified); use -c <commit> to reword the commit message
# u, update-ref <ref> = track a placeholder for the <ref> to be updated
#                       to this position in the new commits. The <ref> is
#                       updated at the end of the rebase
#
# These lines can be re-ordered; they are executed from top to bottom.
#
# If you remove a line here THAT COMMIT WILL BE LOST.
#
# However, if you remove everything, the rebase will be aborted.
#
`[1:]

var squashFileAfter = `
pick bd64edd commit message 1
f 86304e9 commit message 2
f dab68cf commit message 3
f 0c4ebf0 commit message 4

# Rebase a422b17..0c4ebf0 onto a422b17 (4 commands)
#
# Commands:
# p, pick <commit> = use commit
# r, reword <commit> = use commit, but edit the commit message
# e, edit <commit> = use commit, but stop for amending
# s, squash <commit> = use commit, but meld into previous commit
# f, fixup [-C | -c] <commit> = like "squash" but keep only the previous
#                    commit's log message, unless -C is used, in which case
#                    keep only this commit's message; -c is same as -C but
#                    opens the editor
# x, exec <command> = run command (the rest of the line) using shell
# b, break = stop here (continue rebase later with 'git rebase --continue')
# d, drop <commit> = remove commit
# l, label <label> = label current HEAD with a name
# t, reset <label> = reset HEAD to a label
# m, merge [-C <commit> | -c <commit>] <label> [# <oneline>]
#         create a merge commit using the original merge commit's
#         message (or the oneline, if no original merge commit was
#         specified); use -c <commit> to reword the commit message
# u, update-ref <ref> = track a placeholder for the <ref> to be updated
#                       to this position in the new commits. The <ref> is
#                       updated at the end of the rebase
#
# These lines can be re-ordered; they are executed from top to bottom.
#
# If you remove a line here THAT COMMIT WILL BE LOST.
#
# However, if you remove everything, the rebase will be aborted.
#

`[1:]
