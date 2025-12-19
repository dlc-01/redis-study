package response

type ErrorString struct {
	Message string
}

func (ErrorString) isResponse() {}
