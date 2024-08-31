package git

import "github.com/barrettj12/jit/common/types"

type ResetArgs struct {
	Mode   ResetMode
	Branch types.LocalBranch // branch/commit to reset to
}

func Reset(opts ResetArgs) error {
	args := []string{"reset"}
	if opts.Mode != ResetModeUnspecified {
		args = append(args, opts.Mode.String())
	}
	args = append(args, string(opts.Branch))

	_, err := internalExec(internalExecArgs{
		args: args,
	})
	return err
}

type ResetMode string

const (
	ResetModeUnspecified ResetMode = ""

	HardReset ResetMode = "--hard"
)

func (m ResetMode) String() string {
	return string(m)
}
