package path

import "path/filepath"

// Dir represents a directory where Git commands can be run.
type Dir interface {
	Path() string
}

func Path(d Dir) string {
	if d == nil {
		return ""
	}
	return d.Path()
}

// CurrentDir is a default path representing the current directory.
const CurrentDir = currentDir("")

type currentDir string

func (d currentDir) Path() string {
	return ""
}

// JitDir represents the path to the Jit dir.
type JitDir string

func (d JitDir) Path() string {
	return string(d)
}

// Repo represents the path to the root of a Git repository.
type Repo string

func RepoPath(jitDir JitDir, user, repo string) Repo {
	return Repo(filepath.Join(jitDir.Path(), user, repo))
}

func (p Repo) Path() string {
	return string(p)
}

func (p Repo) JitDir() JitDir {
	return JitDir(filepath.Dir(filepath.Dir(string(p))))
}

func (p Repo) Owner() string {
	return filepath.Base(filepath.Dir(string(p)))
}

func (p Repo) RepoName() string {
	return filepath.Base(string(p))
}

// GitFolder represents the path to the .git folder of a repository.
type GitFolder string

func GitFolderPath(path Repo) GitFolder {
	return GitFolder(filepath.Join(path.Path(), ".git"))
}

func (f GitFolder) Path() string {
	return string(f)
}

func (f GitFolder) RepoPath() Repo {
	return Repo(filepath.Dir(string(f)))
}

// Worktree represents the path to a worktree, relative to the repo root.
type Worktree string

func WorktreePath(repo Repo, worktreeName string) Worktree {
	return Worktree(filepath.Join(repo.Path(), worktreeName))
}

func (p Worktree) Path() string {
	return string(p)
}

func (p Worktree) RepoPath() Repo {
	return Repo(filepath.Dir(string(p)))
}

func (p Worktree) WorktreeName() string {
	return filepath.Base(string(p))
}
