package git

func RemoteExists(dir, remote string) (bool, error) {
	_, err := internalExec(internalExecArgs{
		args: []string{"remote", "get-url", remote},
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

func Fetch(dir, remote, branch string) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"fetch", remote, branch},
		dir:  dir,
	})
	return err
}

func AddRemote(remoteName, url string) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"remote", "add", remoteName, url},
	})
	return err
}
