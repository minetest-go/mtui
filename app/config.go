package app

import (
	"bytes"
	"fmt"
	"mtui/minetestconfig"
	"os"
	"path"
	"sync"
)

var mtconfig_mutex = sync.RWMutex{}

func (a *App) ReadMTConfig(sts minetestconfig.SettingTypes) (minetestconfig.Settings, error) {
	mtconfig_mutex.RLock()
	defer mtconfig_mutex.RUnlock()

	data, err := os.ReadFile(a.Config.MinetestConfig)
	if err != nil {
		return nil, fmt.Errorf("error reading config from '%s': %v", a.Config.MinetestConfig, err)
	}

	s := minetestconfig.Settings{}
	err = s.Read(bytes.NewReader(data), sts)
	return s, err
}

func (a *App) WriteMTConfig(cfg minetestconfig.Settings, sts minetestconfig.SettingTypes) error {
	mtconfig_mutex.Lock()
	defer mtconfig_mutex.Unlock()

	f, err := os.OpenFile(a.Config.MinetestConfig, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("could not open minetest config file '%s': %v", a.Config.MinetestConfig, err)
	}
	defer f.Close()

	err = cfg.Write(f, sts)
	if err != nil {
		return fmt.Errorf("could not write minetest config file '%s': %v", a.Config.MinetestConfig, err)
	}

	return nil
}

func (a *App) GetSettingTypes() (minetestconfig.SettingTypes, error) {
	modst, err := minetestconfig.GetAllSettingTypes(path.Join(a.Config.WorldDir, "worldmods"))
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes for worldmods dir: %v", err)
	}

	gamest, err := minetestconfig.GetAllSettingTypes(path.Join(a.Config.WorldDir, "game/mods"))
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes for game/mods dir: %v", err)
	}

	serversettings, err := minetestconfig.GetServerSettingTypes()
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes: %v", err)
	}

	sts := minetestconfig.SettingTypes{}
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
