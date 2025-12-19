package response

type SimpleString struct {
	Value string
}

func (SimpleString) isResponse() {}

type BulkString struct {
	Value *string
}

func (BulkString) isResponse() {}

type ErrorString struct {
	Message string
}

func (ErrorString) isResponse() {}

type Integer struct {
	Value int
}

func (Integer) isResponse() {}

type Array struct {
	Values []Response
}

func (Array) isResponse() {}
