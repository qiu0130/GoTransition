package transitions

import (
	"fmt"
	"log"
)

type Machine struct {
	name                  string
	states                map[string]*State
	events                map[string]*Event
	currentState          *State
	transitions           *Transition
	sendEvent             bool
	ignoreInvalidTriggers bool
	beforeStateChange     []HandleFunc
	afterStateChange      []HandleFunc
	prepareEvent          []HandleFunc
	finalizeEvent         []HandleFunc
}

// new machine
func NewMachine(name, initial string, states []State, trans []Transition, sendEvent, ignoreInvalidTriggers bool,
	beforeStateChange, afterStateChange, prepareEvent, finalizeEvent []HandleFunc) *Machine {

	m := new(Machine)
	if name != "" {
		m.name = "Machine <" + name + ">"
	}
	m.currentState = NewState(initial, ignoreInvalidTriggers, nil, nil)

	m.states = make(map[string]*State, len(states))
	for _, state := range states {
		m.states[state.name] = NewState(state.name, state.ignoreInvalidTriggers, state.onEnter, state.onExit)
	}
	m.events = make(map[string]*Event, len(trans))
	for _, tran := range trans {
		m.events[tran.name] = NewEvent(tran.name, m)
		t := NewTransition(tran.name, tran.source, tran.destination, tran.condition, tran.unless, tran.before, tran.after, tran.prepare)
		m.events[tran.name].addTransition(t)
	}
	m.beforeStateChange = beforeStateChange
	m.afterStateChange = afterStateChange
	m.prepareEvent = prepareEvent
	m.finalizeEvent = finalizeEvent
	m.sendEvent = sendEvent

	return m
}

// machine trigger -> event trigger -> transition execute -> state onEnter and onExit
func (m *Machine) Trigger(name string, args ...interface{}) (*State, error) {

	event, ok := m.events[name]
	if !ok {
		return nil, fmt.Errorf("trigger name: <%s> not found on events", name)
	}
	err := event.trigger(m.currentState.name)
	if err != nil {
		return nil, err
	}
	return m.currentState, nil
}

func (m *Machine) getState(name string) (error, *State) {
	if state, ok := m.states[name]; ok {
		log.Printf("states %v get name %v state %v\n", state, m.states, name)
		return nil, state
	}
	return fmt.Errorf("state '%s' is not registered state", name), nil
}

func (m *Machine) getCurrentState() *State {
	return m.currentState
}

func (m *Machine) setState(name string) error {
	err, state := m.getState(name)
	if err != nil {
		return err
	}
	// update current state
	m.currentState = state
	return nil
}

func (m *Machine) callback(handle HandleFunc, eventData *EventData) error {
	if m.sendEvent {
		handle(eventData)
	}
	return nil
}
