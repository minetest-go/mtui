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
)

type ContentDBModHandler struct{}

func (h *ContentDBModHandler) Create(ctx *HandlerContext, mod *types.Mod) error {
	dir := getDir(ctx.WorldDir, mod)

	cli := cdb.New()
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
		releases, err := cli.GetReleases(mod.Author, mod.Name)
		if err != nil {
			return fmt.Errorf("could not fetch releases: %v", err)
		}
		if len(releases) == 0 {
			return fmt.Errorf("no releases for package '%s/%s'", mod.Author, mod.Name)
		}
		release = releases[0]
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

	// zip structure:
	// mods: "modname/init.lua"
	// games: "mods/modname/init.lua"

	for _, f := range z.File {
		if strings.HasSuffix(f.Name, "/") || strings.HasPrefix(f.Name, "/") {
			// can't do anything with those
			continue
		}

		// the target file to extract to
		var fullpath string
		switch mod.ModType {
		case types.ModTypeMod:
			fullpath = path.Join(path.Dir(dir), f.Name)
		case types.ModTypeGame:
			fullpath = path.Join(dir, f.Name)
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
		if err != nil {
			return fmt.Errorf("could not copy data to '%s': %v", f.Name, err)
		}
	}

	return ctx.Repo.Create(mod)
}

func (h *ContentDBModHandler) Status(ctx *HandlerContext, mod *types.Mod) (*ModStatus, error) {
	return nil, fmt.Errorf("not implemented")
}

func (h *ContentDBModHandler) Update(ctx *HandlerContext, mod *types.Mod, version string) error {
	return fmt.Errorf("not implemented")
}

func (h *ContentDBModHandler) Remove(ctx *HandlerContext, mod *types.Mod) error {
	return fmt.Errorf("not implemented")
}
