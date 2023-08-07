package modmanager

import (
	"mtui/types"
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
	default:
		return ""
	}
}
