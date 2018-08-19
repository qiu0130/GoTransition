package transitions

import (
	"fmt"
	"log"
)

type State struct {
	name                  string
	ignoreInvalidTriggers bool
	onEnter               []HandleFunc
	onExit                []HandleFunc
}

func NewState(name string, ignoreInvalidTriggers bool, onEnter, onExit []HandleFunc) *State {
	return &State{
		name: name,
		ignoreInvalidTriggers: ignoreInvalidTriggers,
		onEnter:               onEnter,
		onExit:                onExit,
	}
}

func (state *State) enter(eventData *EventData) error {
	log.Printf("%s entering state %s, processing callbacks...\n", eventData.machine.name, state.name)
	for _, handle := range state.onEnter {
		eventData.machine.callback(handle, eventData)
	}
	log.Printf("%s entered state %s\n", eventData.machine.name, state.name)
	return nil
}

func (state *State) exit(eventData *EventData) error {
	log.Printf("%s exiting state %s, processiong callbacks...\n", eventData.machine.name, state.name)
	for _, handle := range state.onExit {
		eventData.machine.callback(handle, eventData)
	}
	log.Printf("%s exited state %s\n", eventData.machine.name, state.name)
	return nil
}

func (state *State) addCallback(trigger string, handle HandleFunc) error {

	if trigger == "enter" {
		state.onEnter = append(state.onEnter, handle)
	} else if trigger == "exit" {
		state.onExit = append(state.onExit, handle)
	} else {
		return fmt.Errorf("%s trigger is invalid, only enter or exit", trigger)
	}
	return nil
}
