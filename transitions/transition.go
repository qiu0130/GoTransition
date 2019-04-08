package transitions

import (
	"fmt"
)

type Condition struct {
	handle ConditionFunc
	target bool
}

func (cd *Condition) check(ed *EventData) error {
	// send err event to upstream
	if ed.machine.sendEvent {
		executed := cd.handle(ed)
		if executed == cd.target {
			return nil
		}
		return fmt.Errorf("exected result not equal target")
	}
	return fmt.Errorf("unsupported sendEvent")
}

type Transition struct {
	name        string
	source      string
	destination string
	conditions  []Condition
	prepare     HandleFunc
	before      HandleFunc
	after       HandleFunc

	condition, unless ConditionFunc
}

func NewTransition(name, source, destination string, condition, unless ConditionFunc, prepare, before, after HandleFunc) *Transition {

	tr := new(Transition)
	tr.name = name
	tr.source = source
	tr.destination = destination

	var c []Condition
	if condition != nil {
		c = append(c, Condition{handle: condition, target: true})
	}
	if unless != nil {
		c = append(c, Condition{handle: unless, target: false})
	}

	tr.conditions = c
	tr.prepare = prepare
	tr.before = before
	tr.after = after
	return tr
}

// transition execute
func (tr *Transition) execute(ed *EventData) error {

	Info("%s initiating transition from state <%s> to state <%s>\n", ed.machine.name, tr.source, tr.destination)
	if err := ed.machine.callback(tr.prepare, ed); err != nil {
		Error("failed to prepare func bind callback on transition executing")
		return err
	}
	for _, cond := range tr.conditions {
		err := cond.check(ed)
		if err != nil {
			Error("%s transition condition failed: %v does not return %v\n", ed.machine.name, cond.handle, cond.target)
			return err
		}
	}

	beforeFunc := append(ed.machine.beforeStateChange, tr.before)
	for _, f := range beforeFunc {
		if err := ed.machine.callback(f, ed); err != nil {
			Error("failed to beforeFunc bind callback on transition executing")
			return err
		}
	}

	if err := tr.changeState(ed); err != nil {
		return err
	}

	var afterFunc []HandleFunc
	if tr.after != nil {
		afterFunc = append(afterFunc, tr.after)
	}
	if ed.machine.afterStateChange != nil {
		afterFunc = append(afterFunc, ed.machine.afterStateChange...)
	}
	for _, f := range afterFunc {
		if err := ed.machine.callback(f, ed); err != nil {
			Error("failed to afterFunc bind callback on transition executing")
			return err
		}
	}
	return nil
}

// transition changing state
func (tr *Transition) changeState(ed *EventData) error {
	err, state := ed.machine.getState(tr.source)
	if err != nil {
		return err
	}
	if err := state.exit(ed); err != nil {
		return err
	}
	if err := ed.machine.setState(tr.destination); err != nil {
		return err
	}
	if err := ed.update(ed.state.name); err != nil {
		return err
	}

	err, state = ed.machine.getState(tr.destination)
	if err != nil {
		return err
	}
	if err := state.enter(ed); err != nil {
		return err
	}
	return nil
}

func (tr *Transition) addCallback(trigger string, handle HandleFunc) error {
	switch trigger {
	case "prepare":
		tr.prepare = handle
	case "before":
		tr.before = handle
	case "after":
		tr.after = handle
	default:
		return fmt.Errorf("%s trigger is invalid , only 'prepare', 'before', 'after'", trigger)
	}
	return nil
}
