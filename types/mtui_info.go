package types

type PrivInfo struct {
	Description        string `json:"description"`
	GiveToSingleplayer bool   `json:"give_to_singleplayer"`
	GiveToAdmin        bool   `json:"give_to_admin"`
}

type ChatcommandInfo struct {
	Params      string          `json:"params"`
	Description string          `json:"description"`
	Privs       map[string]bool `json:"privs"`
}
