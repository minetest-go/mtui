package modmanager

import (
	"fmt"
	"mtui/types"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitModHandler struct{}

func (h *GitModHandler) Create(ctx *HandlerContext, mod *types.Mod) error {
	dir := getDir(ctx.WorldDir, mod)

	// prune dir before re-installing
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("error in initial cleanup: %v", err)
	}

	// clone to target dir
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               mod.URL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ShallowSubmodules: true,
	})
	if err != nil {
		return fmt.Errorf("error while cloning: %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("error switching to worktree: %v", err)
	}

	if mod.Branch != "" {
		// check out branch
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(mod.Branch),
		})

		if err == nil && mod.Version != "" {
			// checkout specified hash
			err = w.Reset(&git.ResetOptions{
				Commit: plumbing.NewHash(mod.Version),
				Mode:   git.HardReset,
			})
		}
	} else {
		// default to master or main
		rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
			Name: "origin",
			URLs: []string{mod.URL},
		})
		refs, err := rem.List(&git.ListOptions{})
		if err != nil {
			return fmt.Errorf("git error: %v", err)
		}
		for _, ref := range refs {
			if ref.Name() == plumbing.ReferenceName("refs/heads/master") {
				mod.Branch = "refs/heads/master"
			} else if ref.Name() == plumbing.ReferenceName("refs/heads/main") {
				mod.Branch = "refs/heads/main"
			}
		}
	}

	if mod.Version != "" {
		// check out specified version (no tracking branch)
		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(mod.Version),
		})
	}

	if err != nil {
		return fmt.Errorf("error in checkout: %v", err)
	}

	ref, err := r.Head()
	if err != nil {
		return fmt.Errorf("error switching head: %v", err)
	}

	mod.Version = ref.Hash().String()
	mod.LatestVersion = mod.Version

	return ctx.Repo.Create(mod)
}

func (h *GitModHandler) Update(ctx *HandlerContext, mod *types.Mod, version string) error {
	dir := getDir(ctx.WorldDir, mod)

	r, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("open error: %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("worktree open error: %v", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName:        "origin",
		Depth:             1,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("pull error: %v", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(version),
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("checkout error: %v", err)
	}

	mod.Version = version
	return ctx.Repo.Update(mod)

}

func (h *GitModHandler) Remove(ctx *HandlerContext, mod *types.Mod) error {
	dir := getDir(ctx.WorldDir, mod)

	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return ctx.Repo.Delete(mod.ID)
}

func (h *GitModHandler) CheckUpdate(ctx *HandlerContext, mod *types.Mod) (bool, error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{mod.URL},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("git error: %v", err)
	}
	for _, ref := range refs {
		if ref.Name() == plumbing.ReferenceName(mod.Branch) {
			mod.LatestVersion = ref.Hash().String()
			return true, nil
		}
	}

	return false, nil
}
