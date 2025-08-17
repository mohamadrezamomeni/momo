package tier

type Update struct {
	IdentifyTierDto
	IsDefault *bool `json:"is_default,omitempty"`
}
