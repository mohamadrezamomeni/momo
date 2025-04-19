package dto

type AddInbound struct {
	Port     string
	Tag      string
	Protocol string
	User     *InboundUser
}

type InboundUser struct {
	Username string
	Level    string
	UUID     string
}
