package main

import (
	"reflect"
	"fmt"
	"log"
	"errors"
)

const (
	Separator = "_"
	WildcardAll = "*"
	WildcardSame = "="

)

type HandlerFunc func(eventData *EventData)
type ConditionFunc func(eventData *EventData) bool

type Machine struct {
	model *Machine
	states map[string]*State
	events map[string][]Transition
	initial State
	transitions *Transition
	sendEvent bool
	autoTransitions bool
	orderedTransitions bool
	ignoreInvalidTriggers bool
	beforeStateChange []HandlerFunc
	afterStateChange []HandlerFunc
	name string
	queued bool
	prepareEvent []HandlerFunc
	finalizeEvent []HandlerFunc
}


func (machine *Machine) getState(state string) *State {
	if state, ok := machine.states[state]; ok {
		panic(fmt.Sprintf("state '%s' is not registered state", state))
	}
	return state
}
func (machine *Machine) setState(state String, model string) {


}

func (machine *Machine) getTriggers(states []string) []string {
	for t, ev := range machine.events {
		for _, state := range states {
			if ev.
		}
	}
}

func (machine *Machine) callback(handle HandlerFunc, eventData *EventData) error {
	if machine.sendEvent {
		handle(eventData)
	}
	return nil
}

type State struct {
	name string
	ignoreInvalidTriggers bool
	onEnter []HandlerFunc
	onExit []HandlerFunc
}

func (state *State) enter(eventData *EventData) error {
	log.Println("%s entering state %s, processing callbacks...", eventData.machine.name, state.name)

	for _, handle := range state.onEnter {
		eventData.machine.callback(handle, eventData)
	}
	log.Println("%s entered state %s", eventData.machine.name, state.name)
	return nil
}
func (state *State) exit(eventData *EventData) error {
	log.Println("%s exiting state %s, processiong callbacks...", eventData.machine.name, state.name)
	for _, handle := range state.onExit {
		eventData.machine.callback(handle, eventData)
	}
	log.Println("%s exited state %s", eventData.machine.name, state.name)
	return nil
}

func (state *State) addCallback(trigger string, handle HandlerFunc) error {

	if trigger == "enter" {
		state.onEnter = append(state.onEnter, handle)
	} else if trigger == "exit" {
		state.onExit = append(state.onExit, handle)
	} else {
		return errors.New(fmt.Sprint("%s trigger is invalid, only enter or exit", trigger))
	}
	return nil
}


type Condition struct {
	handle ConditionFunc
	target bool
}
func (cd *Condition) check(ed *EventData) (bool, error) {
	if ed.machine.sendEvent {
		executed := cd.handle(ed)
		if executed == cd.target {
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("unsupported sendEvent")
}
type Transition struct {
	source string
	destination string
	conditions []Condition
	prepare []HandlerFunc
	before []HandlerFunc
	after []HandlerFunc
}

func NewTransition(source, dest string, conditions, unless []ConditionFunc, prepare, before, after []HandlerFunc) *Transition {

	tr := new(Transition)
	tr.source = source
	tr.destination = dest

	var c []Condition
	for _, cond := range conditions {
		c = append(c, Condition{handle:cond, target:true})
	}
	for _, cond := range unless {
		c = append(c, Condition{handle: cond, target: false})
	}
	tr.conditions = c
	tr.prepare = prepare
	tr.before = before
	tr.after = after

	return tr
}

func (tr *Transition) execute(ed *EventData) (error, bool) {

	log.Println("%S initiating transition from state %s to state %s...", ed.machine.name, tr.source, tr.destination)

	for _, f := range tr.prepare {
		ed.machine.callback(f, ed)
		log.Println("executed callback '%s' before conditions ", f)
	}
	for _, cond := range tr.conditions {
		ok, err := cond.check(ed)
		if err != nil {
			return err, false
		}
		if !ok {
			log.Println("%s transition codition failed: %s() does not return %s. transition halted", ed.machine.name, cond.handle, cond.target)
			return nil, false
		}
	}

	beforeFunc := append(ed.machine.beforeStateChange, tr.before...)
	for _, f := range beforeFunc {
		ed.machine.callback(f, ed)
		log.Println("%s executed callback '%s' before transition", ed.machine.name, f)
	}
	tr.changeState(ed)

	afterFunc := append(tr.after, ed.machine.afterStateChange...)
	for _, f := range afterFunc {
		ed.machine.callback(f, ed)
		log.Println("%s executed callback '%s' after transition", ed.machine.name, f)
	}
	return nil, true
}

func (tr *Transition) changeState(ed *EventData) error {
	ed.machine.getState(tr.source).exit(ed)
	ed.machine.setState(tr.destination, ed.model)
	ed.update(ed.model)
	ed.machine.getState(tr.destination).enter(ed)
	return nil
}
func (tr *Transition) addCallback(trigger string, handle HandlerFunc) error {
	switch trigger {
	case "prepare":
		tr.prepare = append(tr.prepare, handle)
	case "before":
		tr.before = append(tr.before, handle)
	case "after":
		tr.after = append(tr.after, handle)
	default:
		return errors.New(fmt.Sprintf("%s trigger is invalid , only 'prepare', 'before', 'after'", trigger))
	}
	return nil
}
type Event struct {
	name string
	machine *Machine
	transitions map[string][]Transition
}
func NewEvent(name string, machine *Machine) *Event {
	return &Event{
		name: name,
		machine: machine,
		transitions: map[string][]Transition{},
	}
}
func (e *Event) addTransition(tr Transition) error {
	e.transitions[tr.source] = append(e.transitions[tr.source], tr)
	return nil
}
type EventData struct {
	state *State
	event *Event
	machine *Machine
	model *Machine
	args []interface{}
	kwargs map[string][]interface{}
	transition *Transition
	error error
	result bool
}

func NewEventData(state *State, event *Event, machine, model *Machine, args []interface{}, kwargs map[string][]interface{}, tr *Transition, error error, result bool) *EventData {

	return &EventData{
		state: state,
		event: event,
		machine: machine,
		model: model,
		args: args,
		kwargs: kwargs,
		transition: tr,
		error: error,
		result: result,
	}
}
func (ed *EventData) update(model *Machine) error {
	ed.state = ed.machine.getState(model.state)
	return nil
}

func main() {
	// tests := []string{
	//	"aa",
	//	"bbb",
	//}
	tests := string("aaa")

	str, ok := listIfy(tests)
	if ok {
		if str != nil {
			fmt.Println(str)
		} else {
			fmt.Println(tests)
		}
	}

}