package statemachine

//go:generate mockgen -source=./IStateMachine.go -destination=./mock/IStateMachineMock.go -package=mock

type IStateMachine interface {
	NextState(string, string) (string, error)
	CommandToEvent(string) (string, error)
	Validate() (bool, error)
	GetStates() []string
	GetActions() []string
	GetEvents() []string
	GetTransitionMatrix() [][]string
	EventsToState([]string) (string, error)
}
