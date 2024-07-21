package git

import (
	"fmt"
	"github.com/barrettj12/jit/common/types"
	"reflect"
	"testing"
)

func TestPull(t *testing.T) {
	// Patch internal exec function
	var expected internalExecArgs
	realInternalExec := internalExec
	internalExec = func(opts internalExecArgs) (string, error) {
		if !reflect.DeepEqual(opts, expected) {
			t.Fatalf("incorrect args %#v\n", opts)
		}
		return "", nil
	}
	t.Cleanup(func() {
		internalExec = realInternalExec
	})

	tests := []struct {
		description string
		pullArgs    PullArgs
		expected    internalExecArgs
	}{{
		description: "pull with explicit remote",
		pullArgs: PullArgs{
			LocalBranch: "branch",
			RemoteBranch: types.RemoteBranch{
				Remote: "remote",
				Branch: "branch",
			},
		},
		expected: internalExecArgs{
			args:         []string{"pull", "remote", "branch:branch"},
			attachStderr: true,
		},
	}, {
		description: "pull with no remote",
		pullArgs: PullArgs{
			LocalBranch:  "branch",
			RemoteBranch: types.NoRemote,
		},
		expected: internalExecArgs{
			args:         []string{"pull"},
			attachStderr: true,
		},
	}}

	for _, test := range tests {
		fmt.Println(test.description)
		expected = test.expected
		_ = Pull(test.pullArgs)
	}
}
