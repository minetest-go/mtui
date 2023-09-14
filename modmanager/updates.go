package modmanager

import (
	"fmt"
)

func (m *ModManager) CheckUpdates() error {
	mods, err := m.repo.GetAll()
	if err != nil {
		return fmt.Errorf("get all mods failed: %v", err)
	}

	for _, mod := range mods {
		h := m.handlers[mod.SourceType]
		updated, err := h.CheckUpdate(m.handlercontext, mod)
		if err != nil {
			return fmt.Errorf("update check failed for mod '%s': %v", mod.Name, err)
		}
		if updated {
			err = m.repo.Update(mod)
			if err != nil {
				return fmt.Errorf("failed to update mod-data for '%s': %v", mod.Name, err)
			}
		}
	}

	return nil
}
