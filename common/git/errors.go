package git

import "strings"

func IsNoSuchRemoteErr(err error) bool {
	return strings.Contains(err.Error(), "error: No such remote ")
}

func IsNoUpstreamConfiguredErr(err error) bool {
	return strings.Contains(err.Error(), "fatal: no upstream configured for branch ")
}
