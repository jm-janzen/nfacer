package common

import (
	"log"
	"os"
	"path"
)

// Change dir to project root, logging this action
func ChdirWebserverRoot() {
	cwd, err := os.Getwd()
	if err != nil {
		Crash(err, 1)
	}
	if cwd == os.Getenv("GOPATH") {
		targetWd := path.Join(cwd, "src", "github.com", "jm-janzen", "nfacer")

		log.Printf("Changing Working Directory from '%v' to '%v'\n", cwd, targetWd)

		os.Chdir(targetWd)
	}
}

// Log error, and exit using given exit code
func Crash(err error, exitCode int) {
	log.Fatal(err)
	os.Exit(exitCode)
}
