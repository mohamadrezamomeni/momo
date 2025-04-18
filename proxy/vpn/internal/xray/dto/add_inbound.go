package dto

type AddInbound struct {
	Port     string
	Tag      string
	Protocol string
	User     *InboundUser
}

type InboundUser struct {
	Email string
	Level string
	UUID  string
}
