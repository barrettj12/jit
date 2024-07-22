package config

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Read a config value from the repo-specific config file at
//
//	<repo-root>/.jit/config.yaml
//
// The provided type parameter defines the type that the value will be
// unmarshalled to.
func readRepoConfig[T any](key string) (t T, err error) {
	repoBasePath, err := getRepoBasePath()
	if err != nil {
		return t, fmt.Errorf("getting repo base path: %w", err)
	}
	repoConfigFilePath := filepath.Join(repoBasePath.Path(), ".jit/config.yaml")
	return readConfigFile[T](repoConfigFilePath, key)
}

// Read a config value from the global Jit config file at
//
//	$HOME/.jit/config.yaml
//
// The provided type parameter defines the type that the value will be
// unmarshalled to.
func readGlobalConfig[T any](key string) (t T, err error) {
	userHomeDir, err := getUserHomeDir()
	if err != nil {
		return t, fmt.Errorf("getting user home dir: %w", err)
	}
	globalConfigFilePath := filepath.Join(userHomeDir, ".jit/config.yaml")
	return readConfigFile[T](globalConfigFilePath, key)
}

// Read a key from a Jit config file. The format of the config file is a
// uni-level mapping of keys to values. The provided type parameter defines the
// type that the value will be unmarshalled to.
func readConfigFile[T any](filePath, key string) (t T, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return t, fmt.Errorf("reading repo config: %w", err)
	}

	var config map[string]T
	// We will get an error for values that don't match type T, but we should
	// ignore this, because we only care about the value corresponding to the
	// given key.
	err = yaml.NewDecoder(file).Decode(&config)
	if _, ok := config[key]; !ok {
		return t, fmt.Errorf("key %q not defined in config file at %s", key, filePath)
	}
	return config[key], nil
}

func addKeyToConfigFile(filePath, key string, value any) {}

// These methods can be replaced for testing.
var getRepoBasePath = common.RepoBasePath
var getUserHomeDir = os.UserHomeDir
