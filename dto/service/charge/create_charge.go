package charge

type CreateChargeDto struct {
	UserID    string
	PackageID string
	InboundID string
	Detail    string
}
