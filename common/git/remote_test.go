package git

import (
	"github.com/barrettj12/jit/common/testutil"
	"testing"
)

func TestListRemotes(t *testing.T) {
	patchInternalExec(t, func(opts internalExecArgs) (string, error) {
		return remotes, nil
	})

	remoteInfo, err := ListRemotes()
	testutil.CheckErr(t, err)
	testutil.AssertDeepEqual(t, remoteInfo, map[string]*RemoteInfo{
		"SimonRichardson": {
			Name:     "SimonRichardson",
			FetchURL: "https://github.com/SimonRichardson/juju",
			PushURL:  "https://github.com/SimonRichardson/juju",
		},
		"barrettj12": {
			Name:     "barrettj12",
			FetchURL: "https://github.com/barrettj12/juju",
			PushURL:  "https://github.com/barrettj12/juju",
		},
		"cderici": {
			Name:     "cderici",
			FetchURL: "https://github.com/cderici/juju",
			PushURL:  "https://github.com/cderici/juju",
		},
		"juju": {
			Name:     "juju",
			FetchURL: "https://github.com/juju/juju",
			PushURL:  "https://github.com/juju/juju",
		},
		"tlm": {
			Name:     "tlm",
			FetchURL: "https://github.com/tlm/juju",
			PushURL:  "https://github.com/tlm/juju",
		},
		"upstream": {
			Name:     "upstream",
			FetchURL: "https://github.com/juju/juju",
			PushURL:  "https://github.com/juju/juju",
		},
		"wallyworld": {
			Name:     "wallyworld",
			FetchURL: "https://github.com/wallyworld/juju",
			PushURL:  "https://github.com/wallyworld/juju",
		},
	})

}

var remotes = `
SimonRichardson	https://github.com/SimonRichardson/juju (fetch)
SimonRichardson	https://github.com/SimonRichardson/juju (push)
barrettj12	https://github.com/barrettj12/juju (fetch)
barrettj12	https://github.com/barrettj12/juju (push)
cderici	https://github.com/cderici/juju (fetch)
cderici	https://github.com/cderici/juju (push)
juju	https://github.com/juju/juju (fetch)
juju	https://github.com/juju/juju (push)
tlm	https://github.com/tlm/juju (fetch)
tlm	https://github.com/tlm/juju (push)
upstream	https://github.com/juju/juju (fetch)
upstream	https://github.com/juju/juju (push)
wallyworld	https://github.com/wallyworld/juju (fetch)
wallyworld	https://github.com/wallyworld/juju (push)
`[1:]
