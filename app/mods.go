package app

import (
	"fmt"
	"mtui/minetestconfig"
	"mtui/types"
	"strings"
)

func (a *App) AddHTTPMod(modname string, cfg minetestconfig.Settings) {
	http_mods := cfg["secure.http_mods"]
	if http_mods == nil || http_mods.Value == "" {
		// create new
		http_mods = &minetestconfig.Setting{
			Value: modname,
		}
		cfg["secure.http_mods"] = http_mods
	} else {
		// append if not in list
		is_in_http_list := false
		for _, mod := range strings.Split(http_mods.Value, ",") {
			if strings.TrimSpace(mod) == modname {
				is_in_http_list = true
				break
			}
		}
		if !is_in_http_list {
			http_mods.Value += fmt.Sprintf(",%s", modname)
		}
	}
}

func (a *App) CreateMTUIMod() (*types.Mod, error) {
	m, err := a.Repos.ModRepo.GetByName("mtui")
	if err != nil {
		return nil, err
	}
	if m == nil {
		// not installed yet
		m = &types.Mod{
			Name:       "mtui",
			ModType:    types.ModTypeMod,
			SourceType: types.SourceTypeGIT,
			URL:        "https://github.com/minetest-go/mtui_mod.git",
			Branch:     "master",
		}
		err = a.ModManager.Create(m)
		if err != nil {
			return nil, fmt.Errorf("error creating mod: %v", err)
		}
	}

	// settings for mtui key/url
	if a.Config.DockerHostname == "" {
		return m, nil
	}

	for _, fname := range []string{types.FEATURE_DOCKER, types.FEATURE_MINETEST_CONFIG} {
		feature, err := a.Repos.FeatureRepository.GetByName(fname)
		if err != nil {
			return nil, fmt.Errorf("feature get error: %v", err)
		}
		if !feature.Enabled {
			return m, nil
		}
	}

	sts, err := a.GetSettingTypes()
	if err != nil {
		return nil, fmt.Errorf("setting types error: %v", err)
	}

	cfg, err := a.ReadMTConfig(sts)
	if err != nil {
		return nil, fmt.Errorf("read config error: %v", err)
	}

	cfg["mtui.url"] = &minetestconfig.Setting{
		Value: fmt.Sprintf("http://%s:8080", a.Config.DockerHostname),
	}
	cfg["mtui.key"] = &minetestconfig.Setting{
		Value: a.Config.APIKey,
	}

	a.AddHTTPMod("mtui", cfg)

	err = a.WriteMTConfig(cfg, sts)
	if err != nil {
		return nil, fmt.Errorf("write config error: %v", err)
	}

	return m, nil
}

func (a *App) CreateBeerchatMod() (*types.Mod, error) {
	m, err := a.Repos.ModRepo.GetByName("beerchat")
	if err != nil {
		return nil, err
	}

	if m == nil {
		// create mod
		m = &types.Mod{
			Name:       "beerchat",
			ModType:    types.ModTypeMod,
			SourceType: types.SourceTypeGIT,
			URL:        "https://github.com/mt-mods/beerchat.git",
			Branch:     "refs/heads/master",
		}
		err = a.ModManager.Create(m)
		if err != nil {
			return nil, fmt.Errorf("error creating mod: %v", err)
		}
	}

	for _, fname := range []string{types.FEATURE_DOCKER, types.FEATURE_MINETEST_CONFIG} {
		feature, err := a.Repos.FeatureRepository.GetByName(fname)
		if err != nil {
			return nil, fmt.Errorf("feature get error: %v", err)
		}
		if !feature.Enabled {
			return m, nil
		}
	}

	sts, err := a.GetSettingTypes()
	if err != nil {
		return nil, fmt.Errorf("setting types error: %v", err)
	}

	cfg, err := a.ReadMTConfig(sts)
	if err != nil {
		return nil, fmt.Errorf("read config error: %v", err)
	}

	cfg["beerchat.matterbridge_url"] = &minetestconfig.Setting{
		Value: "http://matterbridge:4242",
	}
	cfg["beerchat.matterbridge_token"] = &minetestconfig.Setting{
		Value: "my-token",
	}

	a.AddHTTPMod("beerchat", cfg)

	err = a.WriteMTConfig(cfg, sts)
	if err != nil {
		return nil, fmt.Errorf("write config error: %v", err)
	}

	return m, nil
}

func (a *App) CreateMapserverMod() (*types.Mod, error) {
	m, err := a.Repos.ModRepo.GetByName("mapserver")
	if err != nil {
		return nil, err
	}

	if m == nil {
		// create mod
		m = &types.Mod{
			Name:       "mapserver",
			ModType:    types.ModTypeMod,
			SourceType: types.SourceTypeGIT,
			URL:        "https://github.com/minetest-mapserver/mapserver_mod.git",
			Branch:     "refs/heads/master",
		}
		err = a.ModManager.Create(m)
		if err != nil {
			return nil, fmt.Errorf("error creating mod: %v", err)
		}
	}

	for _, fname := range []string{types.FEATURE_DOCKER, types.FEATURE_MINETEST_CONFIG} {
		feature, err := a.Repos.FeatureRepository.GetByName(fname)
		if err != nil {
			return nil, fmt.Errorf("feature get error: %v", err)
		}
		if !feature.Enabled {
			return m, nil
		}
	}

	sts, err := a.GetSettingTypes()
	if err != nil {
		return nil, fmt.Errorf("setting types error: %v", err)
	}

	cfg, err := a.ReadMTConfig(sts)
	if err != nil {
		return nil, fmt.Errorf("read config error: %v", err)
	}

	cfg["mapserver.url"] = &minetestconfig.Setting{
		Value: "http://mapserver:8080",
	}

	a.AddHTTPMod("mapserver", cfg)

	err = a.WriteMTConfig(cfg, sts)
	if err != nil {
		return nil, fmt.Errorf("write config error: %v", err)
	}

	return m, nil
}
