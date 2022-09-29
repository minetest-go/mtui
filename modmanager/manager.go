package modmanager

import (
	"errors"
	"io/ioutil"
	"mtui/db"
	"mtui/types"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/uuid"
)

type ModManager struct {
	world_dir string
	repo      *db.ModRepository
}

func New(world_dir string, repo *db.ModRepository) *ModManager {
	return &ModManager{
		world_dir: world_dir,
		repo:      repo,
	}
}

func (m *ModManager) scanMod(modname, dir string, modtype types.ModType) (bool, error) {
	gitDir := isDir(path.Join(dir, ".git"))
	if !gitDir {
		if isDir(dir) && modtype != types.ModTypeWorldMods {
			// self managed folder
			return true, m.repo.Create(&types.Mod{
				ID:         uuid.NewString(),
				Name:       modname,
				SourceType: types.SourceTypeManual,
				ModType:    modtype,
			})
		}

		return false, nil
	}

	r, err := git.PlainOpen(dir)
	if err != nil {
		return false, err
	}

	rem, err := r.Remote("origin")
	if err != nil {
		return false, err
	}
	if rem == nil {
		return false, errors.New("no origin found")
	}

	ref, err := r.Head()
	if err != nil {
		return false, err
	}

	mod := &types.Mod{
		ID:         uuid.NewString(),
		Name:       modname,
		ModType:    modtype,
		SourceType: types.SourceTypeGIT,
		URL:        rem.Config().URLs[0],
		Branch:     ref.Name().String(),
		Version:    ref.Hash().String(),
	}

	return true, m.repo.Create(mod)
}

func (m *ModManager) scanDir(dir string, modtype types.ModType) error {
	e, err := exists(dir)
	if err != nil {
		return err
	}
	if !e {
		return nil
	}

	l, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range l {
		if fi.IsDir() {
			_, err := m.scanMod(fi.Name(), path.Join(dir, fi.Name()), modtype)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *ModManager) Scan() error {
	// clear mod list
	err := m.repo.DeleteAll()
	if err != nil {
		return err
	}

	found, err := m.scanMod("worldmods", path.Join(m.world_dir, "worldmods"), types.ModTypeWorldMods)
	if err != nil {
		return err
	}

	if !found {
		// worldmods is not a git directory, scan all containing folders
		err := m.scanDir(path.Join(m.world_dir, "worldmods"), types.ModTypeMod)
		if err != nil {
			return err
		}
	}

	err = m.scanDir(path.Join(m.world_dir, "textures"), types.ModTypeTexturepack)
	if err != nil {
		return err
	}

	_, err = m.scanMod("game", path.Join(m.world_dir, "game"), types.ModTypeGame)
	if err != nil {
		return err
	}

	return nil
}

func (m *ModManager) Mods() ([]*types.Mod, error) {
	return m.repo.GetAll()
}

func (m *ModManager) Mod(id string) (*types.Mod, error) {
	return m.repo.GetByID(id)
}

func (m *ModManager) Create(mod *types.Mod) error {
	if mod.ID == "" {
		mod.ID = uuid.NewString()
	}

	if mod.SourceType == types.SourceTypeGIT {
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
	} else {
		return errors.New("source type not implemented")
	}

	return m.repo.Create(mod)
}

func (m *ModManager) Status(mod *types.Mod) (*ModStatus, error) {
	dir := m.getDir(mod)

	status := &ModStatus{}

	if mod.SourceType == types.SourceTypeGIT {
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

func (m *ModManager) Update(mod *types.Mod, version string) error {
	dir := m.getDir(mod)
	if mod.SourceType == types.SourceTypeGIT {
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
	} else {
		return errors.New("source type not implemented")
	}
}

func (m *ModManager) Remove(mod *types.Mod) error {
	dir := m.getDir(mod)
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return m.repo.Delete(mod.ID)
}
