package main

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

func SubmitOrder(args...interface{}) {}

func PayForOrder(args...interface{}) {}

func MoneyEnterIntoAccount(args...interface{}) {}

func Refund(args...interface{}) {}

func BackgroundProcessRefunding(args...interface{}) {}

func main() {

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
					trigger: SubmitOrder,
					name:   "submit_order",
					source: "initialized",
					dest:   "submitted",
					before: submitOrderBefore,
					after: submitOrderAfter,
				},
				Transitions{
					trigger: PayForOrder,
					name: "pay_for_order",
					source: "submitted",
					dest: "all_money_payed",
					condition: paymentCompleted,
					before: payForOrderBefore,
					after: payForOrderAfter,
				},
				Transitions{
					trigger: MoneyEnterIntoAccount,
					name: "money_enter_into_account",
					source: "all_money_payed",
					dest: "stock_up",
					before: moneyEnterIntoAccountBefore,
					after: moneyEnterIntoAccountAfter,
				},
				Transitions{
					trigger: Refund,
					name: "refund",
					source: "stock_up",
					dest: "refund_marked",
					before: refundBefore,
					after: refundAfter,
				},
				Transitions{
					trigger: BackgroundProcessRefunding,
					name: "background_process_refunding",
					source: "refund_marked",
					dest: "closed",
					before: backgroundProcessRefundBefore,
					after: backgroundProcessRefundAfter,
				},

	}
	machine := NewMachine("Stock", "initialized", states, transitions, false, false, nil, nil, nil, nil)

}