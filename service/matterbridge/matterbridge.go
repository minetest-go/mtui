package matterbridge

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

func parseBridge(name string, tree *toml.Tree) error {

	return nil
}

func ParseConfig(data []byte) (*MatterbridgeConfig, error) {
	cfg := &MatterbridgeConfig{}

	t, err := toml.LoadBytes(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing toml data: %v", err)
	}

	irc_list_raw := t.GetArray("irc")
	if irc_list_raw != nil {
		irc_list := irc_list_raw.(*toml.Tree)
		for _, key := range irc_list.Keys() {
			irc_raw := irc_list.Get(key)
			if irc_raw != nil {
				irc := irc_raw.(*toml.Tree)
				fmt.Printf("Key: %s\n", irc.Get("Server"))
			}
		}
	}

	return cfg, nil
}
