package cmd

import (
	"bytes"
	"fmt"
	"github.com/barrettj12/jit/common/path"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"

	"github.com/barrettj12/jit/common"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a Git repo to GitHub",
	RunE:  Publish,
}

func Publish(cmd *cobra.Command, args []string) error {
	// Run git status to check if repo exists
	_, err := common.ExecGit(path.CurrentDir, "status")
	switch err {
	case nil:
		fmt.Println("Git repo already exists")

	case common.ErrGitNotARepo:
		// git init
		initMsg, err := common.ExecGit(path.CurrentDir, "init")
		if err != nil {
			return fmt.Errorf("could not initialise git repo: %w", err)
		}
		fmt.Println(initMsg)

		// Show files to be added
		files, err := common.ExecGit(path.CurrentDir, "status", "-s")
		if err != nil {
			return err
		}
		fmt.Println("New files to be committed:")
		fmt.Println(files)

		ok, err := confirm("Commit the above files")
		if err != nil {
			return err
		}
		if !ok {
			fmt.Print(`
Please create an initial commit, then re-run 
    jit publish
to upload to GitHub.
`[1:])
			return nil
		}

		// git add .
		_, err = common.ExecGit(path.CurrentDir, "add", ".")
		if err != nil {
			return err
		}

		// git commit -m "Initial commit"
		commitMsg, err := common.Prompt(`Enter commit message [default is "Initial commit"]:`)
		if err != nil {
			return err
		}
		if commitMsg == "" {
			commitMsg = "Initial commit"
		}

		_, err = common.ExecGit(path.CurrentDir, "commit", "-m", commitMsg)
		if err != nil {
			return err
		}

	default:
		return err
	}

	// check if remotes exist
	remotes, err := getRemotes()
	if err != nil {
		return err
	}

	if len(remotes) > 0 {
		fmt.Println("remote already exists, nothing to do")
		return nil
	}

	// Get name for new repo
	repoPath, err := common.RepoBasePath()
	if err != nil {
		return err
	}
	defRepoName := repoPath.RepoName()

	repoName, err := common.Prompt(fmt.Sprintf("Name for new repo [default is %q]", defRepoName))
	if err != nil {
		return err
	}
	if repoName == "" {
		repoName = defRepoName
	}

	// gh repo create <name> --source=. --public -r <ghUser>
	remoteName := common.GitHubUser()
	res := common.Exec(common.ExecArgs{
		Cmd:    "gh",
		Args:   []string{"repo", "create", repoName, "--source=.", "--public", "-r", remoteName},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	})
	if err := res.RunError; err != nil {
		return err
	}

	// Push to remote
	// git push <remote> HEAD -u
	_, err = common.ExecGit(path.CurrentDir, "push", remoteName, "HEAD", "-u")
	return err
}

// TODO: move to common
func getRemotes() ([]string, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd := exec.Command("git", "remote")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Read stderr for error info
		errInfo := stderr.String()
		return nil, fmt.Errorf("%s\n%s", errInfo, err)
	}

	remotes := strings.Split(stdout.String(), "\n")
	if remotes[len(remotes)-1] == "" {
		remotes = remotes[:len(remotes)-1]
	}
	return remotes, nil
}
