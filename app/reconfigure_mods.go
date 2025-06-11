package app

import (
	"fmt"
)

func (a *App) ReconfigureSystemMods() error {
	mods, err := a.Repos.ModRepo.GetAll()
	if err != nil {
		return fmt.Errorf("get all mods: %v", err)
	}

	for _, mod := range mods {
		switch mod.Name {
		case "mtui":
			err = a.ModManager.Remove(mod)
			if err != nil {
				return fmt.Errorf("mtui remove failed: %v", err)
			}
			_, err = a.CreateMTUIMod()
			if err != nil {
				return fmt.Errorf("mtui create failed: %v", err)
			}
		case "mapserver":
			err = a.ModManager.Remove(mod)
			if err != nil {
				return fmt.Errorf("mapserver remove failed: %v", err)
			}
			_, err = a.CreateMapserverMod()
			if err != nil {
				return fmt.Errorf("mapserver create failed: %v", err)
			}
		case "beerchat":
			err = a.ModManager.Remove(mod)
			if err != nil {
				return fmt.Errorf("beerchat remove failed: %v", err)
			}
			_, err = a.CreateBeerchatMod()
			if err != nil {
				return fmt.Errorf("beerchat create failed: %v", err)
			}
		}
	}

	return nil
}
