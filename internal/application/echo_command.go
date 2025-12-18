package application

type EchoCommand struct {
	Value string
}

func (EchoCommand) Name() string {
	return "ECHO"
}
