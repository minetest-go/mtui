package modmanager

import (
	"mtui/types"
	"os"
	"path"
)

func (m *ModManager) getDir(mod *types.Mod) string {
	switch mod.ModType {
	case types.ModTypeGame:
		return path.Join(m.world_dir, "game")
	case types.ModTypeMod:
		return path.Join(m.world_dir, "worldmods", mod.Name)
	case types.ModTypeTexturepack:
		return path.Join(m.world_dir, "textures", mod.Name) // TODO: verify
	case types.ModTypeWorldMods:
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

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.IsDir()
	}
	return false
}
