package transitions

import (
	"fmt"
	"testing"
)

func TestNewMachine(t *testing.T) {

	states := []State{
		{
			name: "initialized",
		},
		{
			name: "submitted",
		},
		{
			name: "payed",
		},
	}
	transitions := []Transition{
		{
			name:        "submit_order",
			source:      "initialized",
			destination: "submitted",
			before: func(ed *EventData) {
				fmt.Println("I'am SubmitOrder before callback")
			},
			after: func(ed *EventData) {
				fmt.Println("I'am SubmitOrder after callback")
			},
		},
		{
			name:        "pay_order",
			source:      "submitted",
			destination: "payed",
			before: func(ed *EventData) {
				fmt.Println("I'am PayOrder before callback")
			},
			after: func(ed *EventData) {
				fmt.Println("I'am PayOrder after callback")

			},
			condition: func(ed *EventData) bool {
				fmt.Println("I'am PayOrder condition callback")
				return true
			},
			unless: func(ed *EventData) bool {
				fmt.Println("I'am PayOrder unless callback")
				return false
			},
		},
	}
	var (
		machine *Machine
		state *State
		err error
	)
	machine = NewMachine("order_service", "initialized", states, transitions, true,
		false, nil, nil, nil, nil)

	state, err = machine.Trigger("submit_order")
	if err != nil {
		t.Fatal(err)
	}
	if state.name != "submitted" {
		t.Errorf("expected: %s, result: %s, %s != %s", "submitted", state.name, "submitted", state.name)
	}

	state, err = machine.Trigger("pay_order")
	if err != nil {
		t.Fatal(err)
	}
	if state.name != "payed" {
		t.Errorf("expected %s, result: %s, %s != %s", "payed", state.name, "payed", state.name)
	}

}
