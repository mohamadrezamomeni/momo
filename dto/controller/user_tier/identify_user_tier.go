package usertier

type IdentifyUserTier struct {
	UserID string `param:"user_id"`
	Tier   string `param:"tier"`
}
