package command

import "mtui/bridge"

const (
	COMMAND_SKINS_SET_PNG bridge.CommandType = "set_png_skin"
)

type SkinsSetPNGRequest struct {
	Playername string `json:"playername"`
	SkinName   string `json:"skin_name"`
	PNG        string `json:"png"`
}
