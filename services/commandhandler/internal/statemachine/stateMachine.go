package statemachine

import (
	utils "command-handler/internal/utils"
	"encoding/json"
	"fmt"

	overlog "github.com/Trendyol/overlog"
	"golang.org/x/exp/slices"
)

type StateMachine struct {
	States             []string   `json:"states"`
	Actions            []string   `json:"actions"`
	TransitionMatrix   [][]string `json:"transitions"`
	Events             []string   `json:"events"`
	EventStateRelation []string   `json:"eventStateRelation"`
}

func LoadStateMachine(jsonfileContent []byte) (IStateMachine, error) {

	class, method := "StateMachine", "LoadStateMachine"

	overlog.MDC().Set("method", method)
	overlog.MDC().Set("class", class)
	
	var sm StateMachine
	err := json.Unmarshal(jsonfileContent, &sm)
	if err != nil {
		return nil, fmt.Errorf("parsing state machine config failed. %v", err)
	}
	return &sm, nil
}

func (s *StateMachine) GetStates() []string {
	utils.AddClassAndMethodToMDC(s)
	return s.States
}

func (s *StateMachine) GetActions() []string {
	utils.AddClassAndMethodToMDC(s)
	return s.Actions
}

func (s *StateMachine) GetEvents() []string {
	utils.AddClassAndMethodToMDC(s)
	return s.Events
}

func (s *StateMachine) GetTransitionMatrix() [][]string {
	utils.AddClassAndMethodToMDC(s)
	return s.TransitionMatrix
}

func (s *StateMachine) GetEventStateRelation() []string {
	utils.AddClassAndMethodToMDC(s)
	return s.EventStateRelation
}

func (s *StateMachine) NextState(sourceState string, action string) (string, error) {
	utils.AddClassAndMethodToMDC(s)
	sourceStateIndex := slices.Index(s.States, sourceState)
	if sourceStateIndex == -1 {
		return "", fmt.Errorf("%v is not a valid source state for the loaded state machine", sourceState)
	}
	actionIndex := slices.Index(s.Actions, action)
	if actionIndex == -1 {
		return "", fmt.Errorf("%v is not a valid action for the loaded state machine", action)
	}
	nextState := s.TransitionMatrix[sourceStateIndex][actionIndex]
	if nextState == "None" {
		return "", fmt.Errorf("%v is not a valid action for state %v in the loaded state machine", action, sourceState)
	}
	return nextState, nil
}

func (s *StateMachine) CommandToEvent(action string) (string, error) {
	utils.AddClassAndMethodToMDC(s)
	actionIndex := slices.Index(s.Actions, action)
	if actionIndex == -1 {
		return "", fmt.Errorf("%v is not a valid action for the loaded state machine", action)
	}
	return s.Events[actionIndex], nil
}

func (s *StateMachine) EventsToState(events []string) (string, error) {
	utils.AddClassAndMethodToMDC(s)
	// Here we choose to use ONLY last event to compute state
	eventIndex := slices.Index(s.Events, events[len(events)-1])
	if eventIndex == -1 {
		return "", fmt.Errorf("%v is not a valid event for the loaded state machine", events[len(events)-1])
	}
	return s.EventStateRelation[eventIndex], nil
}

func (s *StateMachine) Validate() (bool, error) {
	utils.AddClassAndMethodToMDC(s)

	// check number of states
	numberOfStates := len(s.States)
	numberofRows := len(s.TransitionMatrix)

	if numberOfStates != numberofRows {
		return false, fmt.Errorf("the number of states in the transition matrix is not equal to the number of states in the states array")
	}

	//check number of actions
	numberOfActions := len(s.Actions)
	numberOfColumns := len(s.TransitionMatrix[0])
	if numberOfActions != numberOfColumns {
		return false, fmt.Errorf("the number of actions in the transition matrix is not equal to the number of actions in the actions array")
	}

	//check existence of final state
	existFinalState := false
	for i := range s.TransitionMatrix {
		for j, cell := range s.TransitionMatrix[i] {
			if cell != "None" {
				break
			}
			if j == len(s.TransitionMatrix[i])-1 {
				existFinalState = true
			}
		}
		if existFinalState {
			break
		}
	}
	if !existFinalState {
		return false, fmt.Errorf("the transition matrix does not contain any final state")
	}

	//check invalid state
	for i := range s.TransitionMatrix {
		for _, cell := range s.TransitionMatrix[i] {
			if !slices.Contains(append(s.States, "None"), cell) {
				return false, fmt.Errorf("state %s is not contained in the transition matrix and is not a valid state", cell)
			}
		}
	}

	//Check initial state
	existInitialState := false
	for _, state := range append(s.States, "None") {
		if !containInMatrix(s.TransitionMatrix, state) {
			existInitialState = true
		}
		if existInitialState {
			break
		}
	}
	if !existInitialState {
		return false, fmt.Errorf("the transition matrix does not contain any initial state")
	}

	//Check number of events equal to number of actions
	numberOfEvent := len(s.Events)
	if numberOfActions != numberOfEvent {
		return false, fmt.Errorf("the number of actions in the actions array is not equal to the number of events in the events array")
	}

	//Check states in eventStateRelation are valid
	for _, state := range s.EventStateRelation {
		if !slices.Contains(s.States, state) {
			return false, fmt.Errorf("state %v in the EventStateRelation doesn't exist in the state array", state)
		}
	}
	if len(s.EventStateRelation) != len(s.Events) {
		return false, fmt.Errorf("eventStateRelation must have same length of event array")
	}

	return true, nil

}

func containInMatrix(matrix [][]string, str string) bool {

	for i := range matrix {
		if slices.Contains(matrix[i], str) {
			return true
		}
	}
	return false
}
