package app

import (
	"bytes"
	"fmt"
	"mtui/luanticonfig"
	"os"
	"path"
	"sync"
)

var mtconfig_mutex = sync.RWMutex{}

func (a *App) ReadMTConfig(sts luanticonfig.SettingTypes) (luanticonfig.Settings, error) {
	if a.Config.LuantiConfig == "" {
		return nil, fmt.Errorf("environment variable for minetest config not configured")
	}

	mtconfig_mutex.RLock()
	defer mtconfig_mutex.RUnlock()

	fi, _ := os.Stat(a.Config.LuantiConfig)
	if fi == nil {
		// init cfg
		err := os.WriteFile(a.Config.LuantiConfig, []byte{}, 0644)
		if err != nil {
			return nil, fmt.Errorf("error initializing config: %v", err)
		}
	} else if fi.Mode().Perm() == 0 {
		// repair configs created with mode 0000 by earlier versions
		err := os.Chmod(a.Config.LuantiConfig, 0644)
		if err != nil {
			return nil, fmt.Errorf("error fixing config permissions: %v", err)
		}
	}

	data, err := os.ReadFile(a.Config.LuantiConfig)
	if err != nil {
		return nil, fmt.Errorf("error reading config from '%s': %v", a.Config.LuantiConfig, err)
	}

	s := luanticonfig.Settings{}
	err = s.Read(bytes.NewReader(data), sts)
	return s, err
}

func (a *App) WriteMTConfig(cfg luanticonfig.Settings, sts luanticonfig.SettingTypes) error {
	mtconfig_mutex.Lock()
	defer mtconfig_mutex.Unlock()

	f, err := os.OpenFile(a.Config.LuantiConfig, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("could not open minetest config file '%s': %v", a.Config.LuantiConfig, err)
	}
	defer f.Close()

	err = cfg.Write(f, sts)
	if err != nil {
		return fmt.Errorf("could not write minetest config file '%s': %v", a.Config.LuantiConfig, err)
	}

	return nil
}

func (a *App) GetSettingTypes() (luanticonfig.SettingTypes, error) {
	modst, err := luanticonfig.GetAllSettingTypes(path.Join(a.Config.WorldDir, "worldmods"))
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes for worldmods dir: %v", err)
	}

	gamest, err := luanticonfig.GetAllSettingTypes(path.Join(a.Config.WorldDir, "game"))
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes for game/mods dir: %v", err)
	}

	serversettings, err := luanticonfig.GetServerSettingTypes()
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes: %v", err)
	}

	sts := luanticonfig.SettingTypes{}
	for k, s := range modst {
		sts[k] = s
	}
	for k, s := range gamest {
		sts[k] = s
	}
	for k, s := range serversettings {
		sts[k] = s
	}

	return sts, nil
}
