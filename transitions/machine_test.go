package transitions

import "testing"

func TestNewMachine(t *testing.T) {

	states := []State{
		State{
			name: "initialized",
		},
		State{
			name: "submitted",
		},
		State{
			name: "payed",
		},
	}
	transitions := []Transition{
		Transition{
			name:        "submit_order",
			source:      "initialized",
			destination: "submitted",
			before: func(ed *EventData) {

			},
			after: func(ed *EventData) {

			},
		},
		Transition{
			name:        "pay_order",
			source:      "submitted",
			destination: "payed",
			before: func(ed *EventData) {

			},
			after: func(ed *EventData) {

			},
			condition: func(ed *EventData) bool {
				return true
			},
			unless: func(ed *EventData) bool {
				return false
			},
		},
	}

	machine := NewMachine("order_system", "initialized", states, transitions, true,
		false, nil, nil, nil, nil)

	state, err := machine.Trigger("submit_order")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(state)
		state, err := machine.Trigger("pay_order")
		if err != nil {
			t.Error(err)
		} else {
			t.Log(state)
		}
	}

}
