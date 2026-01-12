package value

type Type string

const (
	TypeString Type = "string"
	TypeList   Type = "list"
	TypeStream Type = "stream"
)

type Value interface {
	Type() Type
}
