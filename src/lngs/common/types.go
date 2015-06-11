package lngscommon

type Command struct {
	Cmd  string
	Data interface{}
}

type CommandQueue chan *Command
