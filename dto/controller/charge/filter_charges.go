package charge

type FilterCharges struct {
	UserID    string `query:"user_id"`
	InboundID string `query:"inbound_id"`
	Statuses  string `query:"statuses"`
}
