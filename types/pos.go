package types

import "fmt"

type Pos struct {
	X JsonInt `json:"x"`
	Y JsonInt `json:"y"`
	Z JsonInt `json:"z"`
}

func (p *Pos) String() string {
	return fmt.Sprintf("%d/%d/%d", p.X, p.Y, p.Z)
}
