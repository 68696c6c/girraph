package workflow

type State string

const (
	StateNil        State = ""
	StateTodo       State = "todo"
	StateInProgress State = "in_progress"
	StateComplete   State = "complete"
	StateDelayed    State = "delayed"
	StateCancelled  State = "cancelled"
	StateError      State = "error"
)

func stateInStates(state State, states []State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

func stateIsCompletedOrCancelled(state State) bool {
	return stateInStates(state, []State{StateComplete, StateCancelled})
}
