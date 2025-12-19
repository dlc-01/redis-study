package response

type Integer struct {
	Value int
}

func (Integer) isResponse() {}
