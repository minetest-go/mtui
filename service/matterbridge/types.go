package matterbridge

type Bridge struct {
	Type string `json:"type"` // irc, discord, api
	// common
	RemoteNickFormat string `json:"remote_nick_format"`
	// irc
	Server     string `json:"server"`
	Nick       string `json:"nick"`
	Password   string `json:"password"`
	ColorNicks bool   `json:"color_nicks"`
	JoinDelay  int    `json:"join_delay"`
	// discord
	Token        string   `json:"token"`
	UseUserName  bool     `json:"use_username"`
	AllowMention []string `json:"allow_mention"`
	// api
	BindAddress string `json:"bind_address"`
	Buffer      int    `json:"buffer"`
}

type GatewayInOut struct {
	Account string `json:"account"` // reference to Bridge-name
	Channel string `json:"channel"`
}

type Gateway struct {
	InOuts  []*GatewayInOut `json:"in_outs"`
	Name    string          `json:"name"` // channel
	Enabled bool            `json:"enabled"`
}

type MatterbridgeConfig struct {
	Bridges  map[string]*Bridge
	Gateways []*Gateway
}
