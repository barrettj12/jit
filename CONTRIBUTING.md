# `jit` contribution guide

The main entry point to `jit` is via the Bash script `./jit` in the root directory.

New Jit commands are stored in the `cmd/` directory with a `.sh` extension. So, for example, when you call
```
jit foo bar baz
```
Jit will look for the script `cmd/foo.sh`, and if found, it will run this script with the arguments `bar baz`.

If the script is not found, Jit will call
```
git foo bar baz
```