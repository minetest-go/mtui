package modmanager

import (
	"mtui/types"
	"path"
)

func getDir(world_dir string, mod *types.Mod) string {
	switch mod.ModType {
	case types.ModTypeGame:
		return path.Join(world_dir, "game")
	case types.ModTypeMod:
		return path.Join(world_dir, "worldmods", mod.Name)
	case types.ModTypeTexturepack:
		return path.Join(world_dir, "textures", mod.Name) // TODO: verify
	default:
		return ""
	}
}
