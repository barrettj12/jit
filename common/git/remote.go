package git

import (
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
)

func RemoteExists(dir path.Dir, remote types.RemoteName) (bool, error) {
	_, err := internalExec(internalExecArgs{
		args: []string{"remote", "get-url", string(remote)},
		dir:  dir,
	})
	if err == nil {
		return true, nil
	}
	if IsNoSuchRemoteErr(err) {
		return false, nil
	}
	return false, err
}

func Fetch(dir path.Dir, branch types.RemoteBranch) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"fetch", string(branch.Remote), branch.Branch},
		dir:  dir,
	})
	return err
}

func AddRemote(name types.RemoteName, url url.RemoteRepo) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"remote", "add", string(name), url.URL()},
	})
	return err
}
