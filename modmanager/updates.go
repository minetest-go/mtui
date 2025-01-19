package modmanager

import (
	"fmt"
	"mtui/types"
)

func CheckUpdates(world_dir string, mods []*types.Mod) ([]*types.Mod, error) {
	updated_mods := []*types.Mod{}

	for _, mod := range mods {
		h := modHandlers[mod.SourceType]
		updated, err := h.CheckUpdate(world_dir, mod)
		if err != nil {
			return nil, fmt.Errorf("update check failed for mod '%s': %v", mod.Name, err)
		}
		if updated {
			updated_mods = append(updated_mods, mod)
		}
	}

	return updated_mods, nil
}
