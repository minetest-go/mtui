package modmanager

import (
	"fmt"
	"io"
	"mtui/api/cdb"
	"mtui/types"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var cli = cdb.New()
var cached_cli = cdb.NewCachedClient(cli, time.Hour)

type ContentDBModHandler struct{}

func (h *ContentDBModHandler) installMod(world_dir string, mod *types.Mod, release *cdb.PackageRelease) error {
	dir := getDir(world_dir, mod)

	// download release
	z, err := cli.DownloadZip(release)
	if err != nil {
		return fmt.Errorf("could not download zip: %v", err)
	}

	// remove old dir
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("could not remove dir '%s': %v", dir, err)
	}

	// strip the first directory level in the zipfile
	strip_basedir := false
	if mod.ModType == types.ModTypeMod {
		initlua, _ := z.Open("init.lua")
		if initlua == nil {
			strip_basedir = true
		}
	}
	if mod.ModType == types.ModTypeGame {
		gameconf, _ := z.Open("game.conf")
		if gameconf == nil {
			strip_basedir = true
		}
	}

	for _, f := range z.File {
		if strings.HasSuffix(f.Name, "/") || strings.HasPrefix(f.Name, "/") {
			// can't do anything with those
			continue
		}

		// the target file to extract to
		var fullpath string
		switch mod.ModType {
		case types.ModTypeMod, types.ModTypeTexturepack:
			if strip_basedir {
				fullpath = path.Join(path.Dir(dir), f.Name)
			} else {
				fullpath = path.Join(dir, f.Name)
			}
		case types.ModTypeGame:
			if strip_basedir {
				fullpath = path.Join(dir, strings.TrimPrefix(f.Name, fmt.Sprintf("%s/", mod.Name)))
			} else {
				fullpath = path.Join(dir, f.Name)
			}
		default:
			return fmt.Errorf("mod type not supported: %s", mod.ModType)
		}

		// create basedir if it does not exist
		basedir := path.Dir(fullpath)
		err = os.MkdirAll(basedir, 0777)
		if err != nil {
			return fmt.Errorf("could not create directory '%s': %v", basedir, err)
		}

		target, err := os.OpenFile(fullpath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("could not open target file '%s': %v", fullpath, err)
		}

		r, err := f.Open()
		if err != nil {
			return fmt.Errorf("could not open zip-entry '%s': %v", f.Name, err)
		}

		_, err = io.Copy(target, r)
		r.Close()
		target.Close()
		if err != nil {
			return fmt.Errorf("could not copy data to '%s': %v", f.Name, err)
		}
	}

	return nil
}

func (h *ContentDBModHandler) getLatestRelease(mod *types.Mod) (*cdb.PackageRelease, error) {
	releases, err := cli.GetReleases(mod.Author, mod.Name)
	if err != nil {
		return nil, fmt.Errorf("could not fetch releases: %v", err)
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases for package '%s/%s'", mod.Author, mod.Name)
	}
	return releases[0], nil
}

func (h *ContentDBModHandler) Create(world_dir string, mod *types.Mod) error {
	pkg, err := cli.GetDetails(mod.Author, mod.Name)
	if err != nil {
		return fmt.Errorf("could not fetch details: %v", err)
	}
	if pkg == nil {
		return fmt.Errorf("could not find package '%s/%s'", mod.Author, mod.Name)
	}

	var release *cdb.PackageRelease
	if mod.Version == "" {
		// no version specified, fetch latest
		release, err = h.getLatestRelease(mod)
		if err != nil {
			return fmt.Errorf("could not fetch latest release: %v", err)
		}

		mod.Version = fmt.Sprintf("%d", release.ID)

	} else {
		// version specified, fetch specific release
		version, err := strconv.ParseInt(mod.Version, 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse version: '%s'", mod.Version)
		}

		release, err = cli.GetRelease(mod.Author, mod.Name, int(version))
		if err != nil {
			return fmt.Errorf("could not fetch releases: %v", err)
		}
	}

	err = h.installMod(world_dir, mod, release)
	if err != nil {
		return fmt.Errorf("install error: %v", err)
	}

	mod.LatestVersion = mod.Version
	return nil
}

func (h *ContentDBModHandler) Update(world_dir string, mod *types.Mod, version string) error {

	release_id, err := strconv.Atoi(version)
	if err != nil {
		return fmt.Errorf("could not convert version number '%s': %v", version, err)
	}

	release, err := cli.GetRelease(mod.Author, mod.Name, release_id)
	if err != nil {
		return fmt.Errorf("could not get release %d: %v", release_id, err)
	}

	err = h.installMod(world_dir, mod, release)
	if err != nil {
		return fmt.Errorf("could not update mod: %v", err)
	}

	mod.Version = fmt.Sprintf("%d", release.ID)
	return nil
}

func (h *ContentDBModHandler) Remove(world_dir string, mod *types.Mod) error {
	dir := getDir(world_dir, mod)

	// remove dir
	return os.RemoveAll(dir)
}

func (h *ContentDBModHandler) CheckUpdate(world_dir string, mod *types.Mod) (bool, error) {
	updates, err := cached_cli.GetUpdates()
	if err != nil {
		return false, fmt.Errorf("could not get updates: %v", err)
	}

	v := updates[fmt.Sprintf("%s/%s", mod.Author, mod.Name)]
	if v > 0 {
		mod.LatestVersion = fmt.Sprintf("%d", v)
		return true, nil
	}

	return false, nil
}
