package git

import (
	"fmt"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"strings"
)

func CurrentBranch(dir path.Dir) (types.LocalBranch, error) {
	out, err := internalExec(internalExecArgs{
		args: []string{"rev-parse", "--abbrev-ref", "HEAD"},
		dir:  dir,
	})
	if err != nil {
		return "", err
	}
	return types.LocalBranch(strings.TrimSpace(out)), nil
}

// Create a new branch `name` based on `base`.
func CreateBranch(name, base types.LocalBranch) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"branch", string(name), string(base)},
	})
	return err
}

// Retrieves the push target for the specified branch. You can use branch = ""
// for the current branch.
// A return value of "" means no upstream is set.
func PushTarget(branch types.LocalBranch) (types.RemoteBranch, error) {
	out, err := internalExec(internalExecArgs{
		args: []string{"rev-parse", "--abbrev-ref",
			fmt.Sprintf("%s@{push}", branch)},
	})
	if err == nil {
		return types.ParseRemoteBranch(strings.TrimSpace(out)), nil
	}
	if IsNoUpstreamConfiguredErr(err) {
		return types.NoRemote, nil
	}
	return types.NoRemote, err
}

// Retrieves the pull target for the specified branch. You can use branch = ""
// for the current branch.
// A return value of "" means no upstream is set.
func PullTarget(branch types.LocalBranch) (types.RemoteBranch, error) {
	out, err := internalExec(internalExecArgs{
		args: []string{"rev-parse", "--abbrev-ref",
			fmt.Sprintf("%s@{u}", branch)},
	})
	if err == nil {
		return types.ParseRemoteBranch(strings.TrimSpace(out)), nil
	}
	if IsNoUpstreamConfiguredErr(err) {
		return types.NoRemote, nil
	}
	return types.NoRemote, err
}

type PushArgs struct {
	Branch      types.LocalBranch // local branch to push
	Remote      types.RemoteName  // remote to push to
	SetUpstream bool              // should the upstream be set on a successful push
}

func Push(opts PushArgs) error {
	args := []string{"push"}
	if opts.SetUpstream {
		args = append(args, "-u")
	}
	if opts.Remote != "" {
		args = append(args, string(opts.Remote))
	}
	if opts.Branch != "" {
		args = append(args, string(opts.Branch))
	}

	_, err := internalExec(internalExecArgs{
		args:         args,
		attachStderr: true,
	})
	return err
}

type PullArgs struct {
	LocalBranch  types.LocalBranch
	RemoteBranch types.RemoteBranch
	Dir          path.Dir
}

func Pull(opts PullArgs) error {
	args := []string{"pull"}
	if opts.RemoteBranch.Remote != "" {
		args = append(args, string(opts.RemoteBranch.Remote))

		refspec := fmt.Sprintf("%s:%s", opts.RemoteBranch.Branch, opts.LocalBranch)
		args = append(args, refspec)
	}

	_, err := internalExec(internalExecArgs{
		args:         args,
		dir:          opts.Dir,
		attachStderr: true,
	})
	return err
}

type SetUpstreamArgs struct {
	LocalBranch  types.LocalBranch
	RemoteBranch types.RemoteBranch
	Dir          path.Dir
}

func SetUpstream(opts SetUpstreamArgs) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"branch", "-u",
			opts.RemoteBranch.String(), string(opts.LocalBranch)},
		dir: opts.Dir,
	})
	return err
}

func Switch(branch types.LocalBranch) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"switch", string(branch)},
	})
	return err
}

func DeleteBranch(branch types.LocalBranch, force bool) error {
	args := []string{"branch"}
	if force {
		args = append(args, "-D")
	} else {
		args = append(args, "-d")
	}
	args = append(args, string(branch))

	_, err := internalExec(internalExecArgs{
		args: args,
	})
	return err
}
