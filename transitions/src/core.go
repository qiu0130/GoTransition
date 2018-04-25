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
	VERSION = "version 0.0.1"
)

type HandleFunc func(ed *EventData)

type ConditionFunc func(ed *EventData) bool

type TriggerFunc func(args...interface{})

// state defined
type State struct {
	name string
	ignoreInvalidTriggers bool
	onEnter []HandleFunc
	onExit []HandleFunc
}

// new state
func NewState(name string, ignoreInvalidTriggers bool, onEnter, onExit []HandleFunc) *State {
	return &State{
		name: name,
		ignoreInvalidTriggers: ignoreInvalidTriggers,
		onEnter: onEnter,
		onExit: onExit,
	}
}

// enter state
func (state *State) enter(eventData *EventData) error {
	log.Printf("%s entering state %s, processing callbacks...\n", eventData.machine.name, state.name)

	for _, handle := range state.onEnter {
		eventData.machine.callback(handle, eventData)
	}
	log.Printf("%s entered state %s\n", eventData.machine.name, state.name)
	return nil
}
// exit state
func (state *State) exit(eventData *EventData) error {
	log.Printf("%s exiting state %s, processiong callbacks...\n", eventData.machine.name, state.name)
	for _, handle := range state.onExit {
		eventData.machine.callback(handle, eventData)
	}
	log.Printf("%s exited state %s\n", eventData.machine.name, state.name)
	return nil
}

// state add callback
func (state *State) addCallback(trigger string, handle HandleFunc) error {

	if trigger == "enter" {
		state.onEnter = append(state.onEnter, handle)
	} else if trigger == "exit" {
		state.onExit = append(state.onExit, handle)
	} else {
		return errors.New(fmt.Sprintf("%s trigger is invalid, only enter or exit", trigger))
	}
	return nil
}

// condition defined

type Condition struct {
	handle ConditionFunc
	target bool
}

// condition check
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

// transition defined 
type Transition struct {
	name string
	source string
	dest string
	conditions []Condition
	prepare HandleFunc
	before HandleFunc
	after HandleFunc
}

// copy transition
type Transitions struct {
	name string
	trigger TriggerFunc
	source string
	dest string
	prepare HandleFunc
	before HandleFunc
	after HandleFunc
	condition ConditionFunc
	unless ConditionFunc
}

// new transition
func NewTransition(name, source, dest string, condition, unless ConditionFunc, prepare, before, after HandleFunc) *Transition {

	tr := new(Transition)
	tr.name = name
	tr.source = source
	tr.dest = dest

	var c []Condition
	c = append(c, Condition{handle: condition, target: true})
	c = append(c, Condition{handle: unless, target: false})

	tr.conditions = c
	tr.prepare = prepare
	tr.before = before
	tr.after = after
	return tr
}

// transition execute
func (tr *Transition) execute(ed *EventData) (error, bool) {

	log.Printf("%s initiating transition from state %s to state %s...\n", ed.machine.name, tr.source, tr.dest)
	ed.machine.callback(tr.prepare, ed)
	log.Printf("executed callback '%s' before conditions\n", tr.prepare)

	for _, cond := range tr.conditions {
		ok, err := cond.check(ed)
		if err != nil {
			return err, false
		}
		if !ok {
			log.Printf("%s transition codition failed: %v() does not return %s. transition halted\n", ed.machine.name, cond.handle, cond.target)
			return nil, false
		}
	}

	beforeFunc := append(ed.machine.beforeStateChange, tr.before)
	for _, f := range beforeFunc {
		ed.machine.callback(f, ed)
		log.Printf("%s executed callback '%s' before transition\n", ed.machine.name, f)
	}

	tr.changeState(ed)

	var afterFunc []HandleFunc
	afterFunc = append(afterFunc, tr.after)
	afterFunc = append(afterFunc, ed.machine.afterStateChange...)
	for _, f := range afterFunc {
		ed.machine.callback(f, ed)
		log.Printf("%s executed callback '%s' after transition\n", ed.machine.name, f)
	}
	return nil, true
}

// transition change state
func (tr *Transition) changeState(ed *EventData) error {
	ed.machine.getState(tr.source).exit(ed)
	ed.machine.setState(tr.dest)
	ed.update(ed.state.name)
	ed.machine.getState(tr.dest).enter(ed)
	return nil
}

// transition add callback
//func (tr *Transition) addCallback(trigger string, handle HandleFunc) error {
//	switch trigger {
//	case "prepare":
//		tr.prepare = append(tr.prepare, handle)
//	case "before":
//		tr.before = append(tr.before, handle)
//	case "after":
//		tr.after = append(tr.after, handle)
//	default:
//		return errors.New(fmt.Sprintf("%s trigger is invalid , only 'prepare', 'before', 'after'", trigger))
//	}
//	return nil
//}

// eventData defined
type EventData struct {
	state *State
	event *Event
	machine *Machine
	args []interface{}
	kw map[string]interface{}
	transition *Transition
	error string
	result bool
}

// new event data
func NewEventData(state *State, event *Event, machine *Machine, args []interface{},
	kw map[string]interface{}, tr *Transition, error string, result bool) *EventData {

	return &EventData{
		state: state,
		event: event,
		machine: machine,
		args: args,
		kw: kw,
		transition: tr,
		error: error,
		result: result,
	}
}

// event data update
func (ed *EventData) update(state string) error {
	ed.state = ed.machine.getState(state)
	return nil
}

// event defined
type Event struct {
	name string
	machine *Machine
	transitions map[string][]Transition
}

// new event
func NewEvent(name string, m *Machine) *Event {
	return &Event{
		name: name,
		machine: m,
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
func (e *Event) trigger(name string, args []interface{}, kw map[string]interface{}) (error, bool) {
	state := e.machine.getState(name)
	if _, ok := e.transitions[state.name]; !ok {
		msg := fmt.Sprintf("%s Can't trigger event %s from state %s!", e.machine.name, e.name, state.name)
		if state.ignoreInvalidTriggers {
			log.Fatalln(msg)
			return nil, false
		} else {
			panic(msg)

		}
	}
	eventData := &EventData{
		state: state,
		event: e,
		machine: e.machine,
		args: args,
		kw: kw,
	}
	for _, f := range e.machine.prepareEvent {
		e.machine.callback(f, eventData)
		log.Printf("excuted machine preparation callback '%s' before conditions.\n", f)
	}

	defer func() {

		if err := recover(); err != nil {
			eventData.error = fmt.Sprintf("error: %s", err)
		}

		for _, f := range e.machine.finalizeEvent {
			e.machine.callback(f, eventData)
			log.Printf("excuted machine finalize callback '%s'.\n", f)
		}
	}()

	for _, trans := range e.transitions[eventData.state.name] {
		eventData.transition = &trans
		err, ok := trans.execute(eventData)
		if  err != nil {
			panic("execute error")
		}
		if ok {
			eventData.result = true
			break
		}

	}
	return nil, eventData.result

}

// machine defined
type Machine struct {
	name string
	states map[string]*State
	events map[string]*Event
	initial *State
	transitions *Transition
	sendEvent bool
	ignoreInvalidTriggers bool
	beforeStateChange []HandleFunc
	afterStateChange []HandleFunc
	prepareEvent []HandleFunc
	finalizeEvent []HandleFunc
	stateDynamicMethods []string
}

// actor defined
type TransitionDescription struct {
	name string
	trigger HandleFunc
	source string
	dest string
	before, after, prepare HandleFunc
	condition, unless ConditionFunc
}

// 
func convenienceTrigger(model *Machine, triggerName string, args []interface{}, kw map[string]interface{}) (err error, flag bool){

	fn, _ := reflect.TypeOf(model).MethodByName(triggerName)
	params := make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(args)
	params[1] = reflect.ValueOf(kw)
	fn.Func.Call(params)
	return nil, true
}

// new machine
func NewMachine(name, initial string, states []State, trans []Transitions, sendEvent, ignoreInvalidTriggers bool,
	beforeStateChange, afterStateChange, prepareEvent, finalizeEvent []HandleFunc) *Machine {

			m := new(Machine)
			if name != "" {
				m.name = name + ": "
			}
			m.initial = NewState(initial, ignoreInvalidTriggers, nil, nil)

			for _, state := range states {
				m.states[state.name] = &state
			}
			for _, tran := range trans {
				m.events[tran.name] = NewEvent(tran.name, m)
				t := NewTransition(tran.name, tran.source, tran.dest, tran.unless, tran.condition,
					tran.before, tran.after, tran.prepare)
				m.events[tran.name].addTransition(t)
			}
			m.beforeStateChange = beforeStateChange
			m.afterStateChange = afterStateChange
			m.prepareEvent = prepareEvent
			m.finalizeEvent = finalizeEvent
			m.sendEvent = sendEvent
			return m
}

// get state
func (m *Machine) getState(state string) *State {
	if state, ok := m.states[state]; ok {
		panic(fmt.Sprintf("state '%s' is not registered state", state.name))
	}
	return state
}

func (m *Machine) setState(state string) {
	m.state = m.getState(state)
}

func (machine *Machine) getTriggers(states []string) []string {
	for t, ev := range machine.events {
		for _, state := range states {
			if ev.
		}
	}
}

func (machine *Machine) callback(handle HandleFunc, eventData *EventData) error {
	if machine.sendEvent {
		handle(eventData)
	}
	return nil
}


func main() {



}