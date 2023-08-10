package modmanager

import (
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

	// clone to target dir
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               mod.URL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
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
	} else if mod.Version != "" {
		// check out specified version (no tracking branch)
		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(mod.Version),
		})
	}

	if err != nil {
		return err
	}

	ref, err := r.Head()
	if err != nil {
		return err
	}

	mod.Version = ref.Hash().String()

	return ctx.Repo.Create(mod)
}

func (h *GitModHandler) Status(ctx *HandlerContext, mod *types.Mod) (*ModStatus, error) {
	status := &ModStatus{}

	dir := getDir(ctx.WorldDir, mod)

	r, err := git.PlainOpen(dir)
	if err != nil {
		return status, err
	}

	heah, err := r.Head()
	if err != nil {
		return status, err
	}
	status.CurrentVersion = heah.Hash().String()

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{mod.URL},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return status, err
	}
	for _, ref := range refs {
		if ref.Name() == plumbing.ReferenceName(mod.Branch) {
			status.LatestVersion = ref.Hash().String()
		}
	}

	return status, nil
}

func (h *GitModHandler) Update(ctx *HandlerContext, mod *types.Mod, version string) error {
	dir := getDir(ctx.WorldDir, mod)

	r, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(version),
	})
}

func (h *GitModHandler) Remove(ctx *HandlerContext, mod *types.Mod) error {
	dir := getDir(ctx.WorldDir, mod)

	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return ctx.Repo.Delete(mod.ID)
}
