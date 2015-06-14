package lngscmdque

import (
	. "lngs/common"
)

var (
	commandQueues = map[string]CommandQueue{}
)

func GetCommandQueue(name string) CommandQueue {
	cmdque := commandQueues[name]
	if cmdque == nil {
		cmdque = make(CommandQueue)
		commandQueues[name] = cmdque
	}
	return cmdque
}

func PostCommandQueue(name string, cmd *Command) {
	cmdque := GetCommandQueue(name)
	cmdque <- cmd
}
