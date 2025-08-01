package git

import (
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
	"strings"
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

type RemoteInfo struct {
	Name     types.RemoteName
	FetchURL url.GitHubRepo
	PushURL  url.GitHubRepo
}

func ListRemotes() (map[string]*RemoteInfo, error) {
	out, err := internalExec(internalExecArgs{
		args: []string{"remote", "-v"},
	})
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	remotes := map[string]*RemoteInfo{}
	for _, line := range lines {
		split := strings.Split(line, "\t")
		name := split[0]
		if _, ok := remotes[name]; !ok {
			remotes[name] = &RemoteInfo{
				Name: types.RemoteName(name),
			}
		}

		urlInfo := split[1]
		if fetchURL, ok := strings.CutSuffix(urlInfo, " (fetch)"); ok {
			remotes[name].FetchURL = url.GitHubRepo(fetchURL)
		}
		if pushURL, ok := strings.CutSuffix(urlInfo, " (push)"); ok {
			remotes[name].PushURL = url.GitHubRepo(pushURL)
		}
	}
	return remotes, nil
}
