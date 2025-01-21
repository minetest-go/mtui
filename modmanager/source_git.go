package modmanager

import (
	"fmt"
	"mtui/types"
	"os"
	"os/exec"
	"strings"
)

func execGit(workdir string, params []string) ([]byte, error) {
	cmd := exec.Command("git", params...)
	cmd.Dir = workdir
	return cmd.Output()
}

type GitModHandler struct{}

func (h *GitModHandler) Create(world_dir string, mod *types.Mod) error {
	dir := getDir(world_dir, mod)

	// prune dir before re-installing
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("error in initial cleanup: %v", err)
	}

	// clone repo with default branch
	result, err := execGit("/", []string{"clone", "--recurse-submodules", mod.URL, dir})
	if err != nil {
		return fmt.Errorf("clone error: %v, '%s'", err, result)
	}

	if mod.Branch == "" {
		// extract default branch
		result, err = execGit(dir, []string{"rev-parse", "--abbrev-ref", "HEAD"})
		if err != nil {
			return fmt.Errorf("rev-parse branch error: %v, '%s'", err, result)
		}
		mod.Branch = strings.TrimSpace(string(result))
		mod.Branch = strings.ReplaceAll(mod.Branch, "\n", "")
	}

	if mod.Version != "" {
		result, err = execGit(dir, []string{"checkout", mod.Version})
		if err != nil {
			return fmt.Errorf("checkout error: %v, '%s'", err, result)
		}
	}

	result, err = execGit(dir, []string{"rev-parse", "HEAD"})
	if err != nil {
		return fmt.Errorf("rev-parse error: %v, '%s'", err, result)
	}

	mod.Version = strings.TrimSpace(string(result))
	mod.Version = strings.ReplaceAll(mod.Version, "\n", "")
	mod.LatestVersion = mod.Version

	return nil
}

func (h *GitModHandler) Update(world_dir string, mod *types.Mod, version string) error {
	dir := getDir(world_dir, mod)

	result, err := execGit(dir, []string{"fetch", "--recurse-submodules"})
	if err != nil {
		return fmt.Errorf("fetch error: %v, '%s'", err, result)
	}

	result, err = execGit(dir, []string{"checkout", version, "--recurse-submodules"})
	if err != nil {
		return fmt.Errorf("checkout error: %v, '%s'", err, result)
	}

	mod.Version = version
	return nil

}

func (h *GitModHandler) Remove(world_dir string, mod *types.Mod) error {
	dir := getDir(world_dir, mod)
	return os.RemoveAll(dir)
}

func (h *GitModHandler) CheckUpdate(world_dir string, mod *types.Mod) (bool, error) {
	dir := getDir(world_dir, mod)

	previous_version := mod.LatestVersion
	result, err := execGit(dir, []string{"ls-remote", mod.URL, mod.Branch})
	if err != nil {
		return false, fmt.Errorf("ls-remote error: %v, '%s'", err, result)
	}

	str := strings.ReplaceAll(string(result), "\n", "")
	str = strings.TrimSpace(str)
	parts := strings.Split(str, "\t")
	if len(parts) != 2 {
		return false, fmt.Errorf("can't parse result: '%s'", str)
	}
	if len(parts[0]) != 40 {
		return false, fmt.Errorf("can't parse hash: '%s'", parts[0])
	}

	mod.LatestVersion = parts[0]
	return mod.LatestVersion != previous_version, nil
}
