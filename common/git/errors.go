package git

import (
	"regexp"
	"strings"
)

func IsNoSuchRemoteErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "error: No such remote ")
}

func IsRemoteAlreadyExistsErr(err error) bool {
	if err == nil {
		return false
	}
	match, _ := regexp.MatchString("remote .* already exists", err.Error())
	return match
}

func IsNoUpstreamConfiguredErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "fatal: no upstream configured for branch ")
}

func IsNoSuchBranchErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "invalid reference:")
}
