package transitions

import "fmt"


func submitOrderBefore(ed *EventData) {
	fmt.Printf("%s I'am SubmitOrderBefore\n", ed.String())
}

func submitOrderAfter(ed *EventData) {
	fmt.Printf("%s I'am SubmitOrderAfter\n", ed.String())

}

func payForOrderBefore(ed *EventData) {
	fmt.Printf("%s I'am PayForOrderBefore\n", ed.String())

}

func payForOrderAfter(ed *EventData) {
	fmt.Printf("%s I'am PayForOrderAfter\n", ed.String())

}

func moneyEnterIntoAccountBefore(ed *EventData) {
	fmt.Printf("%s I'am MoneyEnterIntoAccountBefore\n", ed.String())

}

func moneyEnterIntoAccountAfter(ed *EventData) {
	fmt.Printf("%s I'am MoneyEnterIntoAccountAfter\n", ed.String())

}

func paymentCompleted(ed *EventData) bool {
	fmt.Printf("%s I'am PaymentCompleted condition\n", ed.String())
	return true
}

func refundBefore(ed *EventData) {
	fmt.Printf("%s I'am refundedBefore\n", ed.String())
}

func refundAfter(ed *EventData) {
	fmt.Printf("%s I'am refundedAfter\n", ed.String())
}

func backgroundProcessRefundBefore(ed *EventData) {
	fmt.Printf("%s I'am BackgroundProcessRefundBefore\n", ed.String())

}

func backgroundProcessRefundAfter(ed *EventData) {
	fmt.Printf("%s I'am BackgroundProcessRefundAfter\n", ed.String())
}

func next(m *Machine, name string) {
	state, err := m.Trigger(name)
	if err != nil {
		fmt.Println(err)
	}
	if state != nil {
		fmt.Println(state.name)
	}
}

func Example() {

	states := []State{
		{
			name: "initialized",
		},
		{
			name: "submitted",
		},
		{
			name: "all_money_payed",
		},
		{
			name: "portion_money_payed",
		},
		{
			name: "stocked_up",
		},
		{
			name: "time_out",
		},
		{
			name: "aborted",
		},
		{
			name: "refund_marked",
		},
		{
			name: "closed",
		},
	}
	transitions := []Transition{
		{
			name:        "submit_order",
			source:      "initialized",
			destination: "submitted",
			before:      submitOrderBefore,
			after:       submitOrderAfter,
		},
		{
			name:        "pay_for_order",
			source:      "submitted",
			destination: "all_money_payed",
			condition:   paymentCompleted,
			before:      payForOrderBefore,
			after:       payForOrderAfter,
		},
		{
			name:        "money_enter_into_account",
			source:      "all_money_payed",
			destination: "stocked_up",
			before:      moneyEnterIntoAccountBefore,
			after:       moneyEnterIntoAccountAfter,
		},
		{
			name:        "refunding",
			source:      "stocked_up",
			destination: "refund_marked",
			before:      refundBefore,
			after:       refundAfter,
		},
		{
			name:        "background_process_refunding",
			source:      "refund_marked",
			destination: "closed",
			before:      backgroundProcessRefundBefore,
			after:       backgroundProcessRefundAfter,
		},
	}
	Debug = false
	machine := NewMachine("order_service", "initialized", states, transitions, true, true,
		nil, nil, nil, nil)


	next(machine, "submit_order")
	next(machine, "pay_for_order")
	next(machine, "money_enter_into_account")
	next(machine, "refunding")
	next(machine, "background_process_refunding")

	// Output:
	// eventData<submit_order> I'am SubmitOrderBefore
	// eventData<submit_order> I'am SubmitOrderAfter
	// submitted
	// eventData<pay_for_order> I'am PayForOrderBefore
	// eventData<pay_for_order> I'am PaymentCompleted condition
	// eventData<pay_for_order> I'am PayForOrderAfter
	// all_money_payed
	// eventData<money_enter_into_account> I'am MoneyEnterIntoAccountBefore
	// eventData<money_enter_into_account> I'am MoneyEnterIntoAccountAfter
	// stocked_up
	// eventData<refunding> I'am refundedBefore
	// eventData<refunding> I'am refundedAfter
	// refund_marked
	// eventData<background_process_refunding> I'am BackgroundProcessRefundBefore
	// eventData<background_process_refunding> I'am BackgroundProcessRefundAfter
	// closed
}


