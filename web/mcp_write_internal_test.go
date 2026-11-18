package web

import (
	"mtui/luanticonfig"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseBanDuration(t *testing.T) {
	for in, expected := range map[string]time.Duration{
		"90s":  90 * time.Second,
		"30m":  30 * time.Minute,
		"12h":  12 * time.Hour,
		"7d":   7 * 24 * time.Hour,
		"1.5d": 36 * time.Hour,
		"2w":   14 * 24 * time.Hour,
	} {
		d, err := parseBanDuration(in)
		assert.NoError(t, err, in)
		assert.Equal(t, expected, d, in)
	}

	for _, in := range []string{"", "d", "xd", "7", "forever"} {
		_, err := parseBanDuration(in)
		assert.Error(t, err, in)
	}
}

func TestValidateSettingValue(t *testing.T) {
	min, max := 1.0, 10.0
	st := func(typ string) *luanticonfig.SettingType {
		return &luanticonfig.SettingType{Type: typ, Min: &min, Max: &max, Choices: []string{"a", "b"}}
	}

	assert.NoError(t, validateSettingValue(st("string"), "anything"))
	assert.NoError(t, validateSettingValue(st("bool"), "true"))
	assert.Error(t, validateSettingValue(st("bool"), "yes"))
	assert.NoError(t, validateSettingValue(st("int"), "5"))
	assert.Error(t, validateSettingValue(st("int"), "5.5"))
	assert.Error(t, validateSettingValue(st("int"), "11"))
	assert.Error(t, validateSettingValue(st("int"), "0"))
	assert.NoError(t, validateSettingValue(st("float"), "2.5"))
	assert.Error(t, validateSettingValue(st("float"), "abc"))
	assert.NoError(t, validateSettingValue(st("enum"), "b"))
	assert.Error(t, validateSettingValue(st("enum"), "c"))
	assert.Error(t, validateSettingValue(st("v3f"), "(1,2,3)"))
}
