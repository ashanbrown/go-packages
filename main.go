package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	rootPackage := strings.TrimSpace(mustRunCmd("sh", "-c", "go list -e ."))

	updatedDirs := map[string]bool{}
	for _, f := range os.Args[1:] {
		dir := filepath.Dir(f)
		updatedDirs[dir] = true
	}

	allPackages := strings.Fields(mustRunCmd("sh", "-c", "go list ./..."))

	coveredPackageMap := map[string]bool{}

	for _, p := range allPackages {
		dirName := strings.TrimPrefix(p, rootPackage+"/")
		if dirName == rootPackage {
			dirName = "."
		}
		if updatedDirs[dirName] {
			coveredPackageMap[p] = true
		}
	}

	var coveredPackages []string
	for p := range coveredPackageMap {
		coveredPackages = append(coveredPackages, p)
	}
	fmt.Printf("%s\n", strings.Join(coveredPackages, " "))
}

func mustRunCmd(name string, args ...string) string {
	if output, err := runCmd(name, args...); err != nil {
		cmd := strings.Join(append([]string{name}, args...), " ")
		fmt.Fprintf(os.Stderr, "Command failed: %s\nError: %s\nOutput: %s\n", cmd, err, output)
		os.Exit(1)
		return ""
	} else {
		return output
	}
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...) // #nosec
	output, err := cmd.CombinedOutput()
	return string(output), err
}
