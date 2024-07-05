package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
)

// Run a Git server using the "smart" HTTP Git protocol
// Run with the env variable GIT_PROJECT_ROOT pointing to a folder containing
// Git repos you want to serve
func main() {
	reposRoot := os.Getenv("GIT_PROJECT_ROOT")
	if reposRoot == "" {
		log.Fatal("GIT_PROJECT_ROOT env variable not set")
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		log.Fatalf("cannot find git: %v", err)
	}
	log.Printf("using git at %q", gitPath)

	gitHandler := &cgi.Handler{
		Path: gitPath,
		Args: []string{"http-backend"},
		Env: []string{
			fmt.Sprintf("GIT_PROJECT_ROOT=%s", reposRoot),
			"GIT_HTTP_EXPORT_ALL=true",
		},
	}

	log.Fatal(http.ListenAndServe(":8080", gitHandler))
}
