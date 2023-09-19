package internal

import (
	sm "command-handler/internal/statemachine"
	es "command-handler/internal/eventstore"
	smConfig "command-handler/internal/statemachine/statemachineconfig"
	"command-handler/internal/utils"
)

var (
	EventStore         es.IEventStore
	StateMachineConfig smConfig.IStateMachineConfig
	StateMachine       sm.IStateMachine
	Timestamp	       utils.ITimestampManager
	ConfigLoaded       bool
)