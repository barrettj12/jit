package types

import (
	"fmt"
	"strings"
)

// LocalBranch represents the name of a local Git branch.
// It can also be HEAD or a commit SHA.
type LocalBranch string

const HEAD LocalBranch = "HEAD"

// RemoteName represents the name of a Git remote.
type RemoteName string

// RemoteBranch represents a remote branch.
type RemoteBranch struct {
	Remote RemoteName
	Branch string
}

var NoRemote = RemoteBranch{}

// Git formats remote branches as remote/branch
func ParseRemoteBranch(s string) RemoteBranch {
	split := strings.SplitN(s, "/", 2)
	return RemoteBranch{
		Remote: RemoteName(split[0]),
		Branch: split[1],
	}
}

func (b RemoteBranch) String() string {
	return fmt.Sprintf("%s/%s", b.Remote, b.Branch)
}

func (b RemoteBranch) AsLocalBranch() LocalBranch {
	return LocalBranch(fmt.Sprintf("%s/%s", b.Remote, b.Branch))
}
