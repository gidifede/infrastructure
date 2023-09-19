package statemachineconfig

//go:generate mockgen -source=./IStateMachineConfig.go -destination=./mock/IStateMachineConfigMock.go -package=mock_stateMachineConfig

import "context"

type IStateMachineConfig interface {
	Get(context.Context) ([]byte, error)
}
