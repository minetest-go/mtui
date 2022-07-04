package modmanager

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type ModManager struct {
	world_dir string
	mods      []*Mod
}

func New(world_dir string) *ModManager {
	return &ModManager{
		world_dir: world_dir,
		mods:      make([]*Mod, 0),
	}
}

func (m *ModManager) scanDir(dir string, modtype ModType) error {
	l, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range l {
		if fi.IsDir() {
			e, err := exists(path.Join(dir, fi.Name(), ".git"))
			if err != nil {
				return err
			}
			if e {
				r, err := git.PlainOpen(path.Join(dir, fi.Name()))
				if err != nil {
					return err
				}

				rem, err := r.Remote("origin")
				if err != nil {
					return err
				}
				if rem == nil {
					return errors.New("no origin found")
				}

				ref, err := r.Head()
				if err != nil {
					return err
				}

				mod := &Mod{
					Name:       fi.Name(),
					ModType:    ModTypeRegular,
					SourceType: SourceTypeGit,
					URL:        rem.Config().URLs[0],
					Branch:     ref.Name().String(),
					Version:    ref.Hash().String(),
				}
				m.mods = append(m.mods, mod)
			}
		}
	}

	return nil
}

func (m *ModManager) Scan() error {
	// TODO: scan other mod types
	err := m.scanDir(path.Join(m.world_dir, "worldmods"), ModTypeRegular)
	if err != nil {
		return err
	}

	return nil
}

func (m *ModManager) Mods() []*Mod {
	return m.mods
}

func (m *ModManager) Create(mod *Mod) error {
	if mod.SourceType == SourceTypeGit {
		dir := m.getDir(mod)
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

		if mod.Version != "" {
			// check out specified version
			err = w.Checkout(&git.CheckoutOptions{
				Hash: plumbing.NewHash(mod.Version),
			})
		} else {
			// check out branch
			err = w.Checkout(&git.CheckoutOptions{
				Branch: plumbing.ReferenceName(mod.Branch),
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
	} else {
		return errors.New("source type not implemented")
	}

	found := false
	for _, lm := range m.mods {
		if lm == mod {
			found = true
			break
		}
	}

	if !found {
		m.mods = append(m.mods, mod)
	}

	return nil
}

func (m *ModManager) Status(mod *Mod) (*ModStatus, error) {
	dir := m.getDir(mod)

	status := &ModStatus{}

	if mod.SourceType == SourceTypeGit {
		r, err := git.PlainOpen(dir)
		if err != nil {
			return status, err
		}

		h, err := r.Head()
		if err != nil {
			return status, err
		}
		status.CurrentVersion = h.Hash().String()

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
	} else {
		return status, errors.New("source type not implemented")
	}
}

func (m *ModManager) Update(mod *Mod, version string) error {
	dir := m.getDir(mod)
	if mod.SourceType == SourceTypeGit {
		r, err := git.PlainOpen(dir)
		if err != nil {
			return err
		}

		w, err := r.Worktree()
		if err != nil {
			return err
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			return err
		}

		return w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(version),
		})
	} else {
		return errors.New("source type not implemented")
	}
}

func (m *ModManager) Remove(mod *Mod) error {
	dir := m.getDir(mod)
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	new_list := make([]*Mod, 0)
	for _, lm := range m.mods {
		if lm != mod {
			new_list = append(new_list, lm)
		}
	}
	m.mods = new_list
	return nil
}
