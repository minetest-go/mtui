package command

import (
	"mtui/bridge"
	"mtui/types"
)

// ui -> game
const COMMAND_ATM_TRANSFER bridge.CommandType = "atm_transfer"

type ATMTransferRequest struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Amount int    `json:"amount"`
}

type ATMTransferResponse struct {
	Success      bool          `json:"success"`
	ErrorMessage string        `json:"errmsg"`
	SourceAmount types.JsonInt `json:"source_amount"`
	TargetAmount types.JsonInt `json:"target_amount"`
}
