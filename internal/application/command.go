package application

type Command struct {
	Name string
	Args []string
	Raw  []byte
}
