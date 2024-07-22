package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TODO: maybe better to do this in each test to isolate the tests
func TestMain(m *testing.M) {
	// Mock out config files
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(fmt.Sprintf("failed to create testing dir: %v", err))
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("error removing temp dir %q: %v", tmpDir, err)
		}
	}()

	getRepoConfigFilePath = func() (string, error) {
		return filepath.Join(tmpDir, "repoconfig"), nil
	}
	getGlobalConfigFilePath = func() (string, error) {
		return filepath.Join(tmpDir, "globalconfig"), nil
	}

	m.Run()
}
