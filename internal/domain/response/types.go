package response

type SimpleString struct {
	Value string
}

func (SimpleString) isResponse() {}

type BulkString struct {
	Value string
}

func (BulkString) isResponse() {}
