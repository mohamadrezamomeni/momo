package user

type Filter struct {
	Users []*UserSerialize `json:"users"`
}
