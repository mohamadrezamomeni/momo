package core

import (
	"github.com/mohamadrezamomeni/momo/pkg/cache"
)

type State struct {
	Path             string
	ControllerStates map[string]any
}

func NewState(path string) *State {
	return &State{
		Path:             path,
		ControllerStates: make(map[string]any),
	}
}

func SetPath(id string, path string) {
	val, isExist := cache.Get(id)
	if !isExist {
		state := NewState(path)
		cache.Set(id, state)
		return
	}

	state, ok := val.(*State)
	if !ok {
		state := NewState(path)
		cache.Set(id, state)
		return
	}

	state.Path = path
	cache.Set(id, state)
}

func DeleteState(id string) {
	cache.Delete(id)
}

func Get(id string) (*State, bool) {
	val, isExist := cache.Get(id)
	if !isExist {
		return nil, false
	}

	state, ok := val.(*State)
	if !ok {
		return nil, false
	}

	return state, true
}

func SetControllerState(id string, key string, controllerState any) {
	state, isExist := Get(id)
	if !isExist {
		state = NewState("")
	}
	state.ControllerStates[key] = controllerState
	replaceState(id, state)
}

func GetControllerState(id string, key string) any {
	state, isExist := Get(id)
	if !isExist {
		state = NewState("")
		replaceState(id, state)
	}
	return state.ControllerStates[key]
}

func replaceState(id string, state *State) {
	cache.Set(id, state)
}
