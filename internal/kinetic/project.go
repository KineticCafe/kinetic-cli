package kinetic

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type Project struct {
	Name string
	Path string
}

func isInGitRepo(path string) bool {
	err := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree", "--quiet").Run()

	return err == nil
}

func gitRepoRoot(path string) (string, error) {
	rootPath, err := exec.Command("git", "-C", path, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(rootPath)), nil
}

func ResolveProject(path string) (Project, error) {
	if isInGitRepo(path) {
		rootPath, err := gitRepoRoot(path)
		if err != nil {
			return Project{}, err
		}

		return Project{filepath.Base(rootPath), rootPath}, nil
	}

	return Project{filepath.Base(path), path}, nil
}
