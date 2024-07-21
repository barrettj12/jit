package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/barrettj12/jit/common/types"
	"os/exec"
)

// PRInfo represents information about a pull request on GitHub
type PRInfo struct {
	BaseBranch string `json:"baseRefName"` // the branch the PR is targeting
}

// GetPRInfo returns information about a pull request based on the given branch.
func GetPRInfo(branch types.LocalBranch) (PRInfo, error) {
	cmd := exec.Command("gh", "pr", "view", string(branch), "--json", "baseRefName")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	var runErr error
	runErr = cmd.Run() // this error contains the exit code

	// handle errors
	if runErr != nil {
		// Read stderr for error info
		errInfo := stderr.String()
		return PRInfo{}, fmt.Errorf("%s\n%w", errInfo, runErr)
	}

	// Unmarshal response to json
	result := PRInfo{}
	err := json.Unmarshal(stdout.Bytes(), &result)
	if err != nil {
		return PRInfo{}, fmt.Errorf("processing json response: %w", err)
	}
	return result, nil
}
