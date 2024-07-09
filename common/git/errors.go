package git

import (
	"regexp"
	"strings"
)

func IsNoSuchRemoteErr(err error) bool {
	return strings.Contains(err.Error(), "error: No such remote ")
}

func IsRemoteAlreadyExistsErr(err error) bool {
	match, _ := regexp.MatchString("remote .* already exists", err.Error())
	return match
}

func IsNoUpstreamConfiguredErr(err error) bool {
	return strings.Contains(err.Error(), "fatal: no upstream configured for branch ")
}
