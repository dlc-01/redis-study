package response

type BulkString struct {
	Value *string
}

func (BulkString) isResponse() {}
