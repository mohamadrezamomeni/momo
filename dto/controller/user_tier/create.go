package usertier

type Create struct {
	UserID string `param:"user_id"`
	Tier   string `param:"tier"`
}
