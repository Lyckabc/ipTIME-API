package enums

import "fmt"

// WifiFrequency represents a wireless band.
type WifiFrequency int

const (
	F2_4GHZ WifiFrequency = iota + 1
	F5GHZ
)

// BSSTag returns the bsstag used in the service.cgi API (e.g. "2g.1", "5g.1").
func (f WifiFrequency) BSSTag(index int) string {
	if f == F2_4GHZ {
		return fmt.Sprintf("2g.%d", index)
	}
	return fmt.Sprintf("5g.%d", index)
}

// Band returns the band string ("2g" or "5g").
func (f WifiFrequency) Band() string {
	if f == F2_4GHZ {
		return "2g"
	}
	return "5g"
}
