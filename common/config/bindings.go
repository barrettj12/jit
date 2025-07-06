package config

import (
	"errors"
)

// Editor returns the default editor.
func Editor() (string, error) {
	editor, err := readRepoConfig[string]("editor")
	if err != nil {
		editor, err = readGlobalConfig[string]("editor")
		if err != nil {
			return "", errors.New("default editor not set")
		}
	}
	return editor, nil
}
