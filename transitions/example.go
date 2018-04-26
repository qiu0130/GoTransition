package transitions

import "fmt"

func submitOrderBefore(ed *EventData) {}

func submitOrderAfter(ed *EventData) {}

func payForOrderBefore(ed *EventData) {}

func payForOrderAfter(ed *EventData) {}

func moneyEnterIntoAccountBefore(ed *EventData) {}

func moneyEnterIntoAccountAfter(ed *EventData) {}

func paymentCompleted(ed *EventData) bool {return true}

func refundBefore(ed *EventData) {}

func refundAfter(ed *EventData) {}

func backgroundProcessRefundBefore(ed *EventData) {}

func backgroundProcessRefundAfter(ed *EventData) {}


func Example() {

	states := []State{
			State{
				name: "initialized",
			},
			State{
				name: "submitted",
			},
			State{
				name: "all_money_payed",
			},
			State{
				name: "portion_money_payed",

			},
			State{
				name: "stocked_up",
			},
			State{
				name: "time_out",

			},
			State{
				name: "aborted",
			},
			State{
				name: "refund_marked",
			},
			State{
				name: "closed",
			},
	}
	transitions := []Transitions{
				Transitions{
					name:   "submit_order",
					source: "initialized",
					dest:   "submitted",
					before: submitOrderBefore,
					after: submitOrderAfter,
				},
				Transitions{
					name: "pay_for_order",
					source: "submitted",
					dest: "all_money_payed",
					condition: paymentCompleted,
					before: payForOrderBefore,
					after: payForOrderAfter,
				},
				Transitions{
					name: "money_enter_into_account",
					source: "all_money_payed",
					dest: "stock_up",
					before: moneyEnterIntoAccountBefore,
					after: moneyEnterIntoAccountAfter,
				},
				Transitions{
					name: "refund",
					source: "stock_up",
					dest: "refund_marked",
					before: refundBefore,
					after: refundAfter,
				},
				Transitions{
					name: "background_process_refunding",
					source: "refund_marked",
					dest: "closed",
					before: backgroundProcessRefundBefore,
					after: backgroundProcessRefundAfter,
				},

	}
	machine := NewMachine("Stock", "initialized", states, transitions, true, true, nil, nil, nil, nil)
	err := machine.Trigger("submit_order")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success")
	}
}