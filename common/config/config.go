package config

import (
	"errors"
	"fmt"
	"github.com/barrettj12/jit/common"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// Return the default editor
func Editor() (string, error) {
	editor, err := readRepoConfig("editor")
	if err != nil {
		editor, err = readGlobalConfig("editor")
		if err != nil {
			return "", errors.New("default editor not set")
		}
	}
	editorStr, ok := editor.(string)
	if !ok {
		return "", fmt.Errorf("invalid value %v for editor", editor)
	}
	return editorStr, nil
}

// Read a config value from the repo-specific config file at
//
//	<repo-root>/.jit/config.yaml
func readRepoConfig(key string) (any, error) {
	repoConfigFilePath, err := getRepoConfigFilePath()
	if err != nil {
		return nil, err
	}
	return readConfigFile(repoConfigFilePath, key)
}

// Read a config value from the global Jit config file at
//
//	$HOME/.jit/config.yaml
func readGlobalConfig(key string) (any, error) {
	globalConfigFilePath, err := getGlobalConfigFilePath()
	if err != nil {
		return nil, err
	}
	return readConfigFile(globalConfigFilePath, key)
}

// Read a key from a Jit config file. The format of the config file is a series
// of nested maps. The provided key is of the form
//
//	key1.key2.key3...
//
// which is interpreted as "take the value of key1 in the base map, then take
// the value of key2 in that map, etc."
func readConfigFile(filePath, key string) (any, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading repo config: %w", err)
	}

	var config map[string]any
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("decoding repo config: %w", err)
	}

	keyPath := strings.Split(key, ".")
	var value any = config
	for _, nextKey := range keyPath {
		var ok bool
		value, ok = value.(map[string]any)[nextKey]
		if !ok {
			return nil, fmt.Errorf("key %q doesn't exist in config", key)
		}
	}
	return value, nil
}

func addKeyToConfigFile(filePath, key string, value any) {}

// These methods can be replaced for testing.
var getRepoConfigFilePath = func() (string, error) {
	repoBasePath, err := common.RepoBasePath()
	if err != nil {
		return "", fmt.Errorf("getting repo base path: %w", err)
	}
	return filepath.Join(repoBasePath, ".jit/config.yaml"), nil
}

var getGlobalConfigFilePath = func() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting user home dir: %w", err)
	}
	return filepath.Join(homeDir, ".jit/config.yaml"), nil
}
