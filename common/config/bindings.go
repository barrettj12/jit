package config

import (
	"errors"
)

// Editor returns the default editor.
func Editor() (string, error) {
	// editor, err := readRepoConfig[string]("editor")
	// if err != nil {
	editor, err := readGlobalConfig[string]("editor")
	if err != nil {
		return "", errors.New("default editor not set")
	}
	// }
	return editor, nil
}

// EditNewBranches returns whether new branches should be automatically opened
// for editing.
func EditNewBranches() (bool, error) {
	editNew, err := readGlobalConfig[bool]("edit-new")
	if err != nil {
		return false, errors.New("edit-new not set")
	}
	return editNew, nil
}
