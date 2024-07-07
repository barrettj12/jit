package git

import "strings"

func IsNoSuchRemoteErr(err error) bool {
	return strings.Contains(err.Error(), "error: No such remote ")
}
