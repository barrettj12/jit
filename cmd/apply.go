package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/barrettj12/jit/common"
)

// Allow applying patches from GitHub
//
//	jit apply <filename>
//	jit apply <url>
func Apply(args []string) error {
	path, err := common.ReqArg(args, 0, "Filepath/URL to patch:")
	if err != nil {
		return err
	}

	_, errStat := os.Stat(path)
	if errStat != nil {
		// Might be a URL
		resp, errGet := http.Get(path)
		if errGet != nil {
			return fmt.Errorf(`could not resolve path %q:
  %v
  %v`, path, errStat, errGet)
		}
		file, err := os.CreateTemp("", "patch-")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
		path = file.Name()
	}

	return common.Git("apply", path)
}
