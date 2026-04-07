package structs

type MacAuthPolicy string

const (
	MacAuthOff       MacAuthPolicy = "off"
	MacAuthWhitelist MacAuthPolicy = "whitelist"
	MacAuthBlacklist MacAuthPolicy = "blacklist"
)

type MacAuth struct {
	BSSTag  string        `json:"bsstag"`
	Policy  MacAuthPolicy `json:"policy"`
	MacList []string      `json:"maclist"`
}
