package modmanager

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
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

func (m *ModManager) Scan() error {
	return nil //TODO
}

func (m *ModManager) Mods() []*Mod {
	return m.mods
}

func (m *ModManager) getDir(mod *Mod) string {
	switch mod.ModType {
	case ModTypeGame:
		return path.Join(m.world_dir, "game")
	case ModTypeRegular:
		return path.Join(m.world_dir, "worldmods", mod.Name)
	case ModTypeTexturepack:
		return path.Join(m.world_dir, "textures", mod.Name) // TODO: verify
	case ModTypeWorldmods:
		return path.Join(m.world_dir, "worldmods")
	default:
		return ""
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (m *ModManager) Create(mod *Mod) error {
	return nil //TODO
}

func (m *ModManager) Status(mod *Mod) (*ModStatus, error) {
	status := &ModStatus{
		CurrentVersion: "",
		LatestVersion:  "",
	}
	return status, nil //TODO
}

func (m *ModManager) Update(mod *Mod, version string) error {
	return nil //TODO
}

func (m *ModManager) Remove(mod *Mod) error {
	return nil //TODO
}

func (m *ModManager) IsSync(mod *Mod) (bool, error) {
	dir := m.getDir(mod)

	// check directory
	isdir, err := exists(dir)
	if err != nil || !isdir {
		return false, err
	}

	if mod.SourceType == SourceTypeGit {
		// check .git directory
		gitdir := path.Join(dir, ".git")
		isdir, err := exists(gitdir)
		if err != nil || !isdir {
			return false, err
		}

		r, err := git.PlainOpen(gitdir)
		if err != nil {
			return false, err
		}

		ref, err := r.Head()
		if err != nil {
			return false, err
		}
		//https://github.com/go-git/go-git/blob/master/_examples/ls-remote/main.go

		fmt.Printf("Hash: %s, Name: %s\n", ref.Hash(), ref.Name())
	}

	return false, nil
}

func (m *ModManager) Sync(mod *Mod) error {
	return nil
}
