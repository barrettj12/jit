package git

func RemoteExists(dir, remote string) (bool, error) {
	_, err := internalExec(dir, "remote", "get-url", remote)
	if err == nil {
		return true, nil
	}
	if IsNoSuchRemoteErr(err) {
		return false, nil
	}
	return false, err
}

func Fetch(dir, remote, branch string) error {
	_, err := internalExec(dir, "fetch", remote, branch)
	return err
}

func AddRemote(remoteName, url string) error {
	_, err := internalExec("", "remote", "add", remoteName, url)
	return err
}
