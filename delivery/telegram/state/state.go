package core

import (
	"github.com/mohamadrezamomeni/momo/pkg/cache"
)

type State struct {
	telegramID string
	idx        int
	Paths      []string
	data       map[string]string
}

func NewState(telegramID string, paths ...string) *State {
	newState := &State{
		telegramID: telegramID,
		Paths:      paths,
		idx:        0,
		data:       make(map[string]string),
	}
	cache.Set(telegramID, newState)
	return newState
}

func FindState(id string) (*State, bool) {
	val, isExisted := cache.Get(id)
	if !isExisted {
		return nil, false
	}

	state := val.(*State)
	return state, true
}

func ResetState(id string) {
	state, isExist := FindState(id)
	if isExist {
		state.Delete()
	}
}

func (s *State) Next() {
	s.idx += 1
	s.Save()
}

func (s *State) Save() {
	cache.Set(s.telegramID, s)
}

func (s *State) Delete() {
	cache.Delete(s.telegramID)
}

func (s *State) IsRequestCompleted() bool {
	if s.idx >= len(s.Paths) {
		return true
	}
	return false
}

func (s *State) ReleaseState() {
	cache.Delete(s.telegramID)
}

func (s *State) GetPath() string {
	return s.Paths[s.idx]
}

func (s *State) SetData(key string, value string) {
	s.data[key] = value
	s.Save()
}

func (s *State) GetData(key string) (string, bool) {
	val, isExist := s.data[key]
	if isExist {
		return val, true
	}
	return "", false
}
