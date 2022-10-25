package types_test

import (
	"encoding/json"
	"mtui/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyStruct struct {
	Value types.JsonInt `json:"value"`
}

func TestFloatJson(t *testing.T) {
	str := `{"value": 10.0}`
	s := &MyStruct{}

	dec := json.NewDecoder(strings.NewReader(str))
	dec.UseNumber()
	err := dec.Decode(s)

	assert.NoError(t, err)
	assert.Equal(t, types.JsonInt(10), s.Value)

	s.Value++
}
