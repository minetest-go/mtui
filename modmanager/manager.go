package modmanager

import (
	"errors"
	"fmt"
	"mtui/api/cdb"
	"mtui/db"
	"mtui/types"
	"os"
	"strconv"
	"strings"

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

	dir := m.getDir(mod)

	switch mod.SourceType {
	case types.SourceTypeGIT:
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

	case types.SourceTypeCDB:
		parts := strings.Split(mod.Name, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid modname: '%s'", mod.Name)
		}
		author := parts[0]
		package_name := parts[1]

		cli := cdb.New()
		pkg, err := cli.GetDetails(author, package_name)
		if err != nil {
			return fmt.Errorf("could not fetch details: %v", err)
		}
		if pkg == nil {
			return fmt.Errorf("could not find package '%s/%s'", author, package_name)
		}

		var release *cdb.PackageRelease
		if mod.Version == "" {
			// no version specified, fetch latest
			releases, err := cli.GetReleases(author, package_name)
			if err != nil {
				return fmt.Errorf("could not fetch releases: %v", err)
			}
			if len(releases) == 0 {
				return fmt.Errorf("no releases for package '%s/%s'", author, package_name)
			}
			release = releases[0]

		} else {
			// version specified, fetch specific release
			version, err := strconv.ParseInt(mod.Version, 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse version: '%s'", mod.Version)
			}

			release, err = cli.GetRelease(author, package_name, int(version))
			if err != nil {
				return fmt.Errorf("could not fetch releases: %v", err)
			}
		}

		z, err := cli.DownloadZip(release)
		fmt.Print(z)
		//TODO

	default:
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
