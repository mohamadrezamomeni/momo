package event

type CreateEvent struct {
	Name        string
	Data        string
	IsProcessed bool
}
