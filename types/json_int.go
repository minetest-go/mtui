package types

import "encoding/json"

type JsonInt int64

func JsonIntPtr(i int64) *JsonInt {
	var ji JsonInt = JsonInt(i)
	return &ji
}

func (i *JsonInt) UnmarshalJSON(b []byte) error {
	var v float64
	err := json.Unmarshal(b, &v)
	*i = JsonInt(v)
	return err
}
