package types

type MatterbridgeBridgeIRC struct {
	Name             string `json:"name"` // internal
	Server           string `json:"server"`
	Nick             string `json:"nick"`
	Password         string `json:"password"`
	ColorNicks       bool   `json:"color_nicks"`
	JoinDelay        int    `json:"join_delay"`
	RemoteNickFormat string `json:"remote_nick_format"`
}

type MatterbridgeBridgeDiscord struct {
	Name             string   `json:"name"` // internal
	Server           string   `json:"server"`
	Token            string   `json:"token"`
	UseUserName      bool     `json:"use_username"`
	AllowMention     []string `json:"allow_mention"`
	RemoteNickFormat string   `json:"remote_nick_format"`
}

type MatterbridgeBridgeApi struct {
	Name             string `json:"name"` // internal
	BindAddress      string `json:"bind_address"`
	Token            string `json:"token"`
	Buffer           int    `json:"buffer"`
	RemoteNickFormat string `json:"remote_nick_format"`
}

type MatterbridgeGatewayInOut struct {
	Account string `json:"account"`
	Channel string `json:"channel"`
}

type MatterbridgeGateway struct {
	InOuts  []*MatterbridgeGatewayInOut `json:"in_outs"`
	Name    string                      `json:"name"` // channel
	Enabled bool                        `json:"enabled"`
}

type MatterbridgeConfig struct {
	IRCBridges     []*MatterbridgeBridgeIRC
	DiscordBridges []*MatterbridgeBridgeDiscord
	ApiBridges     []*MatterbridgeBridgeApi
	Gateways       []*MatterbridgeGateway
}
