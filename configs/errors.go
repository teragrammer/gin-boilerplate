package configs

type ErrorContent struct {
	Code    string
	Message string
}

type ErrorId struct {
	E1  ErrorContent
	E2  ErrorContent
	E3  ErrorContent
	E4  ErrorContent
	E5  ErrorContent
	E6  ErrorContent
	E7  ErrorContent
	E8  ErrorContent
	E9  ErrorContent
	E10 ErrorContent
}

func Errors() ErrorId {
	return ErrorId{
		E1:  ErrorContent{Code: "e1", Message: "The application token has an incorrect length"},
		E2:  ErrorContent{Code: "e2", Message: "Too many requests"},
		E3:  ErrorContent{Code: "e3", Message: "Binding form failed"},
		E4:  ErrorContent{Code: "e4", Message: "Validation errors were encountered during the process"},
		E5:  ErrorContent{Code: "e5", Message: "The data you've selected is already assigned"},
		E6:  ErrorContent{Code: "e6", Message: "Password is not correctly formatted"},
		E7:  ErrorContent{Code: "e7", Message: "Whoops something went wrong"},
		E8:  ErrorContent{Code: "e8", Message: "Oops something went wrong during setting initialization"},
		E9:  ErrorContent{Code: "e9", Message: "The information you are looking for is not available"},
		E10: ErrorContent{Code: "e10", Message: "The system has detected too many incorrect login attempts"},
	}
}
