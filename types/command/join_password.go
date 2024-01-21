package command

const COMMAND_SET_JOIN_PASSWORD = "set_join_password"

type SetJoinPasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
