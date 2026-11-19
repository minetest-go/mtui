package app

import (
	"fmt"
	"mtui/luanticonfig"
	"mtui/types"
	"os"
	"strings"
)

func (a *App) AddHTTPMod(modname string, cfg luanticonfig.Settings) {
	http_mods := cfg["secure.http_mods"]
	if http_mods == nil || http_mods.Value == "" {
		// create new
		http_mods = &luanticonfig.Setting{
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

// ensureGitMod installs the git mod if missing and, if the luanti-config feature
// is enabled, applies the config changes made by configure to the engine config
func (a *App) ensureGitMod(name, url, branch string, configure func(cfg luanticonfig.Settings) error) (*types.Mod, error) {
	m, err := a.Repos.ModRepo.GetByName(name)
	if err != nil {
		return nil, err
	}

	if m == nil {
		// not installed yet
		m = &types.Mod{
			Name:       name,
			ModType:    types.ModTypeMod,
			SourceType: types.SourceTypeGIT,
			URL:        url,
			Branch:     branch,
		}
		err = a.ModManager.Create(m)
		if err != nil {
			return nil, fmt.Errorf("error creating mod: %v", err)
		}
	}

	feature, err := a.Repos.FeatureRepository.GetByName(types.FEATURE_LUANTI_CONFIG)
	if err != nil {
		return nil, fmt.Errorf("feature get error: %v", err)
	}
	if !feature.Enabled {
		return m, nil
	}

	sts, err := a.GetSettingTypes()
	if err != nil {
		return nil, fmt.Errorf("setting types error: %v", err)
	}

	cfg, err := a.ReadMTConfig(sts)
	if err != nil {
		return nil, fmt.Errorf("read config error: %v", err)
	}

	err = configure(cfg)
	if err != nil {
		return nil, err
	}

	a.AddHTTPMod(name, cfg)

	err = a.WriteMTConfig(cfg, sts)
	if err != nil {
		return nil, fmt.Errorf("write config error: %v", err)
	}

	return m, nil
}

func (a *App) CreateMTUIMod() (*types.Mod, error) {
	return a.ensureGitMod("mtui", "https://github.com/luanti-golang/mtui_mod.git", "master", func(cfg luanticonfig.Settings) error {
		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("could not get hostname: %v", err)
		}
		cfg["mtui.url"] = &luanticonfig.Setting{Value: fmt.Sprintf("http://%s:8080", hostname)}
		cfg["mtui.key"] = &luanticonfig.Setting{Value: a.Config.APIKey}
		return nil
	})
}

func (a *App) CreateBeerchatMod() (*types.Mod, error) {
	return a.ensureGitMod("beerchat", "https://github.com/mt-mods/beerchat.git", "refs/heads/master", func(cfg luanticonfig.Settings) error {
		cfg["beerchat.matterbridge_url"] = &luanticonfig.Setting{Value: "http://matterbridge:4242"}
		cfg["beerchat.matterbridge_token"] = &luanticonfig.Setting{Value: "my-token"}
		return nil
	})
}

func (a *App) CreateMapserverMod() (*types.Mod, error) {
	return a.ensureGitMod("mapserver", "https://github.com/minetest-mapserver/mapserver_mod.git", "refs/heads/master", func(cfg luanticonfig.Settings) error {
		cfg["mapserver.url"] = &luanticonfig.Setting{Value: "http://mapserver:8080"}
		return nil
	})
}

// ProvisionGame installs the game from the given git url as mod "game".
// Skipped if a "game" mod is already installed; a failed/incomplete one is replaced
func (a *App) ProvisionGame(url string) error {
	existing, err := a.Repos.ModRepo.GetByName("game")
	if err != nil {
		return err
	}
	if existing != nil {
		if existing.Status == types.ModStatusInstalled {
			// already installed, nothing to do
			return nil
		}
		// failed or incomplete install, retry
		err = a.ModManager.Remove(existing)
		if err != nil {
			return fmt.Errorf("error removing existing game: %v", err)
		}
	}

	m := &types.Mod{
		Name:       "game",
		ModType:    types.ModTypeGame,
		SourceType: types.SourceTypeGIT,
		URL:        url,
	}
	err = a.ModManager.Create(m)
	if err != nil {
		return fmt.Errorf("error creating game: %v", err)
	}
	return nil
}
