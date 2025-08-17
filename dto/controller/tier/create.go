package tier

type CreateTier struct {
	IdentifyTierDto
	IsDefault *bool `json:"is_default,omitempty"`
}
