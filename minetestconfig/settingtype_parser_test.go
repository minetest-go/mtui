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
