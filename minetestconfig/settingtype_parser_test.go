package minetestconfig_test

import (
	"mtui/minetestconfig"
	"testing"

	"github.com/stretchr/testify/assert"
)

const minetest_setting = `

[General]

[*Advanced]

[**Networking]

#    Prometheus listener address.
#    If Minetest is compiled with ENABLE_PROMETHEUS option enabled,
#    enable metrics listener for Prometheus on that address.
#    Metrics can be fetched on http://127.0.0.1:30000/metrics
prometheus_listener_address (Prometheus listener address) string 127.0.0.1:30000

#    Maximum size of the out chat queue.
#    0 to disable queueing and -1 to make the queue size unlimited.
max_out_chat_queue_size (Maximum size of the out chat queue) int 20 -1 32767

#    Timeout for client to remove unused map data from memory, in seconds.
client_unload_unused_data_timeout (Mapblock unload timeout) float 600.0 0.0
`

func TestSettingParserminetestSettings(t *testing.T) {
	list, err := minetestconfig.ParseSettingTypes([]byte(minetest_setting))
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))

	assert.Equal(t, 3, len(list[0].Category))
	assert.Equal(t, "General", list[0].Category[0])
	assert.Equal(t, "Advanced", list[0].Category[1])
	assert.Equal(t, "Networking", list[0].Category[2])
	assert.Equal(t, "prometheus_listener_address", list[0].Key)
	assert.Equal(t, "Prometheus listener address", list[0].ShortDescription)
	assert.Equal(t, "127.0.0.1:30000", list[0].Default)
	assert.Equal(t, "string", list[0].Type)

	assert.Equal(t, 3, len(list[1].Category))
	assert.Equal(t, "General", list[1].Category[0])
	assert.Equal(t, "Advanced", list[1].Category[1])
	assert.Equal(t, "Networking", list[1].Category[2])
	assert.Equal(t, "max_out_chat_queue_size", list[1].Key)
	assert.Equal(t, "Maximum size of the out chat queue", list[1].ShortDescription)
	assert.Equal(t, "20", list[1].Default)
	assert.Equal(t, "int", list[1].Type)
	assert.Equal(t, float64(-1), list[1].Min)
	assert.Equal(t, float64(32767), list[1].Max)

	assert.Equal(t, 3, len(list[2].Category))
	assert.Equal(t, "General", list[2].Category[0])
	assert.Equal(t, "Advanced", list[2].Category[1])
	assert.Equal(t, "Networking", list[2].Category[2])
	assert.Equal(t, "client_unload_unused_data_timeout", list[2].Key)
	assert.Equal(t, "Mapblock unload timeout", list[2].ShortDescription)
	assert.Equal(t, "float", list[2].Type)
	assert.Equal(t, float64(0.0), list[2].Min)

}

const mod_settings = `
# Allows the wrench to be crafted if either the 'technic' or 'default' mod is installed.
wrench.enable_crafting (Enable crafting recipe) bool true

# The number of times the wrench can be used before breaking. Default 50. Set to 0 for infinite uses.
wrench.tool_uses (Wrench uses) int 50

# Enables compression of item metadata when picking up nodes. Significantly decreases the size of item metadata, at the cost of not being human-readable.
wrench.compress_data (Compress item metadata) bool true
`

func TestSettingParserMod(t *testing.T) {
	list, err := minetestconfig.ParseSettingTypes([]byte(mod_settings))
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))

	assert.Equal(t, 0, len(list[0].Category))
	assert.Equal(t, "wrench.enable_crafting", list[0].Key)
	assert.Equal(t, "Enable crafting recipe", list[0].ShortDescription)
	assert.Equal(t, "bool", list[0].Type)
	assert.Equal(t, "true", list[0].Default)

	assert.Equal(t, 0, len(list[1].Category))
	assert.Equal(t, "wrench.tool_uses", list[1].Key)
	assert.Equal(t, "Wrench uses", list[1].ShortDescription)
	assert.Equal(t, "int", list[1].Type)
	assert.Equal(t, "50", list[1].Default)

	assert.Equal(t, 0, len(list[2].Category))
	assert.Equal(t, "wrench.compress_data", list[2].Key)
	assert.Equal(t, "Compress item metadata", list[2].ShortDescription)
	assert.Equal(t, "bool", list[2].Type)
	assert.Equal(t, "true", list[2].Default)

}

const setting_with_enum = `
[Appearance]
# Specifies how the value indicators (i.e. health, breah, etc.) look. There are 3 styles
# available. You can choose between the default progress-bar-like bars and the good
# old statbars like you know from vanilla Minetest.
# These values are possible:
#   - progress_bar:    A horizontal progress-bar-like bar with a label, showing numerical value
#                      (current, maximum), and an icon. These bars usually convey the most
#                      information. This is the default and recommended value.
#   - statbar_classic: Classic statbar, like in vanilla Minetest. Made out of up to 20
#                      half-symbols. Those bars represent the vague ratio between
#                      the current value and the maximum value. 1 half-symbol stands for
#                      approximately 5% of the maximum value.
#   - statbar_modern:  Like the classic statbar, but also supports background images, this
#                      kind of statbar may be considered to be more user-friendly than the
#                      classic statbar. This bar type closely resembles the mod
#                      “Better HUD” [hud] by BlockMen.
hudbars_bar_type (HUD bars style) enum progress_bar progress_bar,statbar_classic,statbar_modern
`

func TestEnumSetting(t *testing.T) {
	list, err := minetestconfig.ParseSettingTypes([]byte(setting_with_enum))
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	assert.Equal(t, 1, len(list[0].Category))
	assert.Equal(t, "Appearance", list[0].Category[0])
	assert.Equal(t, "hudbars_bar_type", list[0].Key)
	assert.Equal(t, "HUD bars style", list[0].ShortDescription)
	assert.Equal(t, "enum", list[0].Type)
	assert.Equal(t, 3, len(list[0].Choices))
	assert.Equal(t, "progress_bar", list[0].Choices[0])
	assert.Equal(t, "statbar_classic", list[0].Choices[1])
	assert.Equal(t, "statbar_modern", list[0].Choices[2])
}

func TestNoiseParams2D(t *testing.T) {
	list, err := minetestconfig.ParseSettingTypes([]byte("mgfractal_np_seabed (Seabed noise) noise_params_2d -14, 9, (600, 601, 602), 41900, 5, 0.6, 2.0, eased"))
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	assert.Equal(t, "mgfractal_np_seabed", list[0].Key)
	assert.Equal(t, "Seabed noise", list[0].ShortDescription)
	assert.Equal(t, "noise_params_2d", list[0].Type)
	assert.Equal(t, -14.0, list[0].Offset)
	assert.Equal(t, 9.0, list[0].Scale)
	assert.Equal(t, 600.0, list[0].SpreadX)
	assert.Equal(t, 601.0, list[0].SpreadY)
	assert.Equal(t, 602.0, list[0].SpreadZ)
	assert.Equal(t, "41900", list[0].Seed)
	assert.Equal(t, 5.0, list[0].Octaves)
	assert.Equal(t, 0.6, list[0].Persistence)
	assert.Equal(t, 2.0, list[0].Lacunarity)
	assert.Equal(t, 1, len(list[0].DefaultMGFlags))
	assert.Equal(t, "eased", list[0].DefaultMGFlags[0])
}

func TestTypeV3F(t *testing.T) {
	list, err := minetestconfig.ParseSettingTypes([]byte("mgfractal_scale (Scale) v3f (4096.0, 1024.0, 2048.0)"))
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	assert.Equal(t, "Scale", list[0].ShortDescription)
	assert.Equal(t, "mgfractal_scale", list[0].Key)
	assert.Equal(t, "v3f", list[0].Type)
	assert.Equal(t, 4096.0, list[0].X)
	assert.Equal(t, 1024.0, list[0].Y)
	assert.Equal(t, 2048.0, list[0].Z)
}

const nested_settings = `
[D1]
[*D2]
[**D3]

[T1]

[*T2]

x2 (desc2) string

[**D4]

[**T3]

x3 (desc3) string
`

func TestNestedSettings(t *testing.T) {
	list, err := minetestconfig.ParseSettingTypes([]byte(nested_settings))
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 2, len(list))

	assert.Equal(t, 2, len(list[0].Category))
	assert.Equal(t, "x2", list[0].Key)
	assert.Equal(t, "T1", list[0].Category[0])
	assert.Equal(t, "T2", list[0].Category[1])

	assert.Equal(t, 3, len(list[1].Category))
	assert.Equal(t, "x3", list[1].Key)
	assert.Equal(t, "T1", list[1].Category[0])
	assert.Equal(t, "T2", list[1].Category[1])
	assert.Equal(t, "T3", list[1].Category[2])
}

func TestGetServerSettingTypes(t *testing.T) {
	ss, err := minetestconfig.GetServerSettingTypes()
	assert.NoError(t, err)
	assert.NotNil(t, ss)

	var mgv7_spflags *minetestconfig.SettingType
	for _, s := range ss {
		if s.Key == "mgv7_spflags" {
			mgv7_spflags = s
			break
		}
	}

	assert.NotNil(t, mgv7_spflags)
	assert.Equal(t, 2, len(mgv7_spflags.Category))
	assert.Equal(t, "Mapgen", mgv7_spflags.Category[0])
	assert.Equal(t, "Mapgen V7", mgv7_spflags.Category[1])
}
