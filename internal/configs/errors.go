package configs

type ErrorContent struct {
	Code    string
	Message string
}

type ErrorId struct {
	E1 ErrorContent
	E2 ErrorContent
	E3 ErrorContent
	E4 ErrorContent
}

func Errors() ErrorId {
	return ErrorId{
		E1: ErrorContent{Code: "e1", Message: "The application token has an incorrect length"},
		E2: ErrorContent{Code: "e2", Message: "Too many requests"},
		E3: ErrorContent{Code: "e3", Message: "Binding form failed"},
		E4: ErrorContent{Code: "e4", Message: "Validation errors were encountered during the process"},
	}
}
