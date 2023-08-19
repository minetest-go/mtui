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
	sts, err := minetestconfig.ParseSettingTypes([]byte(minetest_setting))
	assert.NoError(t, err)
	assert.NotNil(t, sts)
	assert.Equal(t, 3, len(sts))

	e := sts["prometheus_listener_address"]
	assert.Equal(t, 3, len(e.Category))
	assert.Equal(t, "General", e.Category[0])
	assert.Equal(t, "Advanced", e.Category[1])
	assert.Equal(t, "Networking", e.Category[2])
	assert.Equal(t, "prometheus_listener_address", e.Key)
	assert.Equal(t, "Prometheus listener address", e.ShortDescription)
	assert.Equal(t, "127.0.0.1:30000", e.Default.Value)
	assert.Equal(t, "string", e.Type)

	e = sts["max_out_chat_queue_size"]
	assert.Equal(t, 3, len(e.Category))
	assert.Equal(t, "General", e.Category[0])
	assert.Equal(t, "Advanced", e.Category[1])
	assert.Equal(t, "Networking", e.Category[2])
	assert.Equal(t, "max_out_chat_queue_size", e.Key)
	assert.Equal(t, "Maximum size of the out chat queue", e.ShortDescription)
	assert.Equal(t, "20", e.Default.Value)
	assert.Equal(t, "int", e.Type)
	assert.Equal(t, float64(-1), e.Min)
	assert.Equal(t, float64(32767), e.Max)

	e = sts["client_unload_unused_data_timeout"]
	assert.Equal(t, 3, len(e.Category))
	assert.Equal(t, "General", e.Category[0])
	assert.Equal(t, "Advanced", e.Category[1])
	assert.Equal(t, "Networking", e.Category[2])
	assert.Equal(t, "client_unload_unused_data_timeout", e.Key)
	assert.Equal(t, "Mapblock unload timeout", e.ShortDescription)
	assert.Equal(t, "float", e.Type)
	assert.Equal(t, float64(0.0), e.Min)

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
	sts, err := minetestconfig.ParseSettingTypes([]byte(mod_settings))
	assert.NoError(t, err)
	assert.NotNil(t, sts)
	assert.Equal(t, 3, len(sts))

	e := sts["wrench.enable_crafting"]
	assert.Equal(t, 0, len(e.Category))
	assert.Equal(t, "wrench.enable_crafting", e.Key)
	assert.Equal(t, "Enable crafting recipe", e.ShortDescription)
	assert.Equal(t, "bool", e.Type)
	assert.Equal(t, "true", e.Default.Value)

	e = sts["wrench.tool_uses"]
	assert.Equal(t, 0, len(e.Category))
	assert.Equal(t, "wrench.tool_uses", e.Key)
	assert.Equal(t, "Wrench uses", e.ShortDescription)
	assert.Equal(t, "int", e.Type)
	assert.Equal(t, "50", e.Default.Value)

	e = sts["wrench.compress_data"]
	assert.Equal(t, 0, len(e.Category))
	assert.Equal(t, "wrench.compress_data", e.Key)
	assert.Equal(t, "Compress item metadata", e.ShortDescription)
	assert.Equal(t, "bool", e.Type)
	assert.Equal(t, "true", e.Default.Value)

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
	sts, err := minetestconfig.ParseSettingTypes([]byte(setting_with_enum))
	assert.NoError(t, err)
	assert.NotNil(t, sts)
	assert.Equal(t, 1, len(sts))

	e := sts["hudbars_bar_type"]
	assert.Equal(t, 1, len(e.Category))
	assert.Equal(t, "Appearance", e.Category[0])
	assert.Equal(t, "hudbars_bar_type", e.Key)
	assert.Equal(t, "HUD bars style", e.ShortDescription)
	assert.Equal(t, "enum", e.Type)
	assert.Equal(t, 3, len(e.Choices))
	assert.Equal(t, "progress_bar", e.Choices[0])
	assert.Equal(t, "statbar_classic", e.Choices[1])
	assert.Equal(t, "statbar_modern", e.Choices[2])
}

func TestNoiseParams2D(t *testing.T) {
	sts, err := minetestconfig.ParseSettingTypes([]byte("mgfractal_np_seabed (Seabed noise) noise_params_2d -14, 9, (600, 601, 602), 41900, 5, 0.6, 2.0, eased"))
	assert.NoError(t, err)
	assert.NotNil(t, sts)
	assert.Equal(t, 1, len(sts))

	e := sts["mgfractal_np_seabed"]
	assert.Equal(t, "mgfractal_np_seabed", e.Key)
	assert.Equal(t, "Seabed noise", e.ShortDescription)
	assert.Equal(t, "noise_params_2d", e.Type)
	assert.Equal(t, -14.0, e.Default.Offset)
	assert.Equal(t, 9.0, e.Default.Scale)
	assert.Equal(t, 600.0, e.Default.SpreadX)
	assert.Equal(t, 601.0, e.Default.SpreadY)
	assert.Equal(t, 602.0, e.Default.SpreadZ)
	assert.Equal(t, "41900", e.Default.Seed)
	assert.Equal(t, 5.0, e.Default.Octaves)
	assert.Equal(t, 0.6, e.Default.Persistence)
	assert.Equal(t, 2.0, e.Default.Lacunarity)
	assert.Equal(t, 1, len(e.DefaultMGFlags))
	assert.Equal(t, "eased", e.DefaultMGFlags[0])
}

func TestTypeV3F(t *testing.T) {
	sts, err := minetestconfig.ParseSettingTypes([]byte("mgfractal_scale (Scale) v3f (4096.0, 1024.0, 2048.0)"))
	assert.NoError(t, err)
	assert.NotNil(t, sts)
	assert.Equal(t, 1, len(sts))

	e := sts["mgfractal_scale"]
	assert.Equal(t, "Scale", e.ShortDescription)
	assert.Equal(t, "mgfractal_scale", e.Key)
	assert.Equal(t, "v3f", e.Type)
	assert.Equal(t, 4096.0, e.Default.X)
	assert.Equal(t, 1024.0, e.Default.Y)
	assert.Equal(t, 2048.0, e.Default.Z)
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
	sts, err := minetestconfig.ParseSettingTypes([]byte(nested_settings))
	assert.NoError(t, err)
	assert.NotNil(t, sts)
	assert.Equal(t, 2, len(sts))

	e := sts["x2"]
	assert.Equal(t, 2, len(e.Category))
	assert.Equal(t, "x2", e.Key)
	assert.Equal(t, "T1", e.Category[0])
	assert.Equal(t, "T2", e.Category[1])

	e = sts["x3"]
	assert.Equal(t, 3, len(e.Category))
	assert.Equal(t, "x3", e.Key)
	assert.Equal(t, "T1", e.Category[0])
	assert.Equal(t, "T2", e.Category[1])
	assert.Equal(t, "T3", e.Category[2])
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
