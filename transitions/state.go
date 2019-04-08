package transitions

import (
	"fmt"
)

type State struct {
	name                  string
	ignoreInvalidTriggers bool
	onEnter               []HandleFunc
	onExit                []HandleFunc
}

func NewState(name string, ignoreInvalidTriggers bool, onEnter, onExit []HandleFunc) *State {
	return &State{
		name:                  name,
		ignoreInvalidTriggers: ignoreInvalidTriggers,
		onEnter:               onEnter,
		onExit:                onExit,
	}
}

func (state *State) enter(eventData *EventData) error {

	Info("%s entering state %s, processing callbacks...\n", eventData.machine.name, state.name)
	for _, handle := range state.onEnter {
		if err := eventData.machine.callback(handle, eventData); err != nil {
			Error("failed to eventData bind callback on enter")
			return err
		}
	}
	return nil
}

func (state *State) exit(eventData *EventData) error {
	Info("%s exiting state %s, processing callbacks...\n", eventData.machine.name, state.name)
	for _, handle := range state.onExit {
		if err := eventData.machine.callback(handle, eventData); err != nil {
			Error("failed to eventData bind callback on exist")
			return err
		}
	}
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
