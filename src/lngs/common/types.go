package lngscommon

type Command struct {
	EntityId string
	Command  string
	Data     interface{}
}

type CommandQueue chan *Command
