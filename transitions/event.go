package transitions

import (
	"fmt"
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

func (ed *EventData) String() string {
	return fmt.Sprintf("eventData<%s>", ed.event.name)
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

func (e *Event) addTransition(tr *Transition) error {
	e.transitions[tr.source] = append(e.transitions[tr.source], *tr)
	return nil
}

func (e *Event) addCallback(trigger string, handle HandleFunc) error {
	var values []Transition
	for _, v := range e.transitions {
		values = append(values, v...)
	}
	for _, v := range values {
		if err := v.addCallback(trigger, handle); err != nil {
			Error(err.Error())
			return err
		}
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
			Info(err.Error())
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
		Info("executed machine preparation callback %v before conditions\n", f)
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
		if err := e.machine.callback(f, eventData); err != nil {
			Error("failed to finalizeEvent bind callback on trigger")
		}
		Info("executed machine finalize callback %v\n", f)
	}
	return nil

}
