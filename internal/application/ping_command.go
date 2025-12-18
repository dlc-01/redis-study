package application

type PingCommand struct{}

func (PingCommand) Name() string {
	return "PING"
}
