package app


type Output struct {
	Success bool
	Error error
	IsUnexpected bool
	Message string
	Data interface{}
}

func NewSuccessOutput() Output{
	return Output{
		Success: true,
	}
}

func NewErrorOutput(err error) Output {
	return Output{
		Error: err,
	}
}

func NewUnexpectedErrorOutput() Output {
	return Output{
		IsUnexpected: true,
	}
}