package testutil

import (
	"os/exec"
	"strings"
	"testing"
)

// RunCommand runs a command in the given dir, and asserts that it returns
// an exit code of zero.
func RunCommand(t *testing.T, dir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf(`error running command %q: %v
output: %s`, strings.Join(append([]string{name}, args...), " "), err, string(out))
	}
}

// CheckErr asserts that the given error is nil.
func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func AssertEqual[T comparable](t *testing.T, obtained, expected T) {
	if obtained != expected {
		t.Fatalf("obtained %#v, expected %#v", obtained, expected)
	}
}

func AssertNotEqual[T comparable](t *testing.T, obtained, expected T) {
	if obtained == expected {
		t.Fatalf("obtained and expected are equal = %#v", expected)
	}
}
