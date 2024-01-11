package command

const COMMAND_GET_CONTROLS_METADATA = "get_controls_metadata"

type ControlsMetadata struct {
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Label       string  `json:"label"`
	Description string  `json:"description"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	Step        float64 `json:"step"`
}

type GetControlsMetadataResponse map[string]*ControlsMetadata

const COMMAND_GET_CONTROLS_VALUES = "get_controls_values"

type GetControlsValuesResponse map[string]any

const COMMAND_SET_CONTROL = "set_control"

type SetControlRequest struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
