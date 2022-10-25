package types

import "encoding/json"

type JsonInt int64

func (i *JsonInt) UnmarshalJSON(b []byte) error {
	var v float64
	err := json.Unmarshal(b, &v)
	*i = JsonInt(v)
	return err
}
