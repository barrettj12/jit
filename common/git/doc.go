// Package git provides nice Go bindings for various Git commands.
//
// The functions in this package should not do much more than format
// inputs/outputs and run the git command.
//
// Ideally, code in the cmd package should never call git directly - instead,
// it should always go through a function in this package.
package git
