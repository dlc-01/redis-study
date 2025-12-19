package commands

type PingCommand struct{}

func (PingCommand) Name() string {
	return "PING"
}
