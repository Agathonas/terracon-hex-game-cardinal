// component/hex.go
package component

type Hex struct {
	Q int `json:"q"` // Column (also known as the x coordinate)
	R int `json:"r"` // Row (also known as the y coordinate)
	S int `json:"s"` // The third coordinate (can be calculated as -Q-R)
}

// Name returns the name of the component, satisfying the Component interface.
func (h Hex) Name() string {
	return "Hex"
}

// NewHex creates a new Hex component with the provided coordinates.
func NewHex(q, r int) Hex {
	return Hex{Q: q, R: r, S: -q - r}
}
