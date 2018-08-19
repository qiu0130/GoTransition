package transitions

import (
	"fmt"
	"log"
)

// current state, event, machine and transition
type EventData struct {
	state      *State
	event      *Event
	machine    *Machine
	transition *Transition
	err        error
	args       []interface{}
}

func NewEventData(state *State, event *Event, machine *Machine, tr *Transition, err error, args []interface{}) *EventData {
	return &EventData{
		state:      state,
		event:      event,
		machine:    machine,
		transition: tr,
		err:        err,
		args:       args,
	}
}

// update current machine state
func (ed *EventData) update(name string) error {
	err, state := ed.machine.getState(name)
	if err != nil {
		return err
	}
	ed.state = state
	return nil
}

// event
type Event struct {
	name        string
	machine     *Machine
	transitions map[string][]Transition
}

func NewEvent(name string, m *Machine) *Event {
	return &Event{
		name:        name,
		machine:     m,
		transitions: map[string][]Transition{},
	}
}

// add transition
func (e *Event) addTransition(tr *Transition) error {
	e.transitions[tr.source] = append(e.transitions[tr.source], *tr)
	return nil
}

// add callback
func (e *Event) addCallback(trigger string, handle HandleFunc) error {
	var values []Transition
	for _, v := range e.transitions {
		values = append(values, v...)
	}
	for _, v := range values {
		v.addCallback(trigger, handle)
	}
	return nil
}

// event trigger
func (e *Event) trigger(name string, args ...interface{}) error {
	err, state := e.machine.getState(name)
	if err != nil {
		return err
	}
	if _, ok := e.transitions[state.name]; !ok {
		err = fmt.Errorf("%s can't trigger event %s from state %s", e.machine.name, e.name, state.name)
		// ignore invalid trigger err
		if state.ignoreInvalidTriggers {
			return err
		}
		panic(err)
	}
	eventData := &EventData{
		state:   state,
		event:   e,
		machine: e.machine,
		args:    args,
	}
	for _, f := range e.machine.prepareEvent {
		err := e.machine.callback(f, eventData)
		if err != nil {
			return err
		}
		log.Printf("excuted machine preparation callback %q before conditions\n", f)
	}

	for _, trans := range e.transitions[eventData.state.name] {
		eventData.transition = &trans
		err := trans.execute(eventData)
		if err != nil {
			eventData.err = err
			return err
		}
		return nil
	}
	for _, f := range e.machine.finalizeEvent {
		e.machine.callback(f, eventData)
		log.Printf("excuted machine finalize callback %q\n", f)
	}
	return nil

}
