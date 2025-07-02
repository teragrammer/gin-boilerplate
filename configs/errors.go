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
	E11 ErrorContent
	E12 ErrorContent
	E13 ErrorContent
	E14 ErrorContent
	E15 ErrorContent
	E16 ErrorContent
	E17 ErrorContent
	E18 ErrorContent
	E19 ErrorContent
	E20 ErrorContent
	E21 ErrorContent
	E22 ErrorContent
	E23 ErrorContent
	E24 ErrorContent
	E25 ErrorContent
	E26 ErrorContent
	E27 ErrorContent
	E28 ErrorContent
	E29 ErrorContent
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
		E11: ErrorContent{Code: "e11", Message: "Your account has been temporarily locked due to multiple login attempts"},
		E12: ErrorContent{Code: "e12", Message: "Token format does not meet requirements"},
		E13: ErrorContent{Code: "e13", Message: "The token has an incorrect length"},
		E14: ErrorContent{Code: "e14", Message: "User information is incorrect"},
		E15: ErrorContent{Code: "e15", Message: "The token detail entered is invalid"},
		E16: ErrorContent{Code: "e16", Message: "The validity of the token is not valid"},
		E17: ErrorContent{Code: "e17", Message: "Token timestamp could not be processed"},
		E18: ErrorContent{Code: "e18", Message: "The token is no longer valid"},
		E19: ErrorContent{Code: "e19", Message: "No OTP is necessary for this process"},
		E20: ErrorContent{Code: "e20", Message: "There is an issue with your email configuration"},
		E21: ErrorContent{Code: "e21", Message: "Resending the OTP is currently not possible"},
		E22: ErrorContent{Code: "e22", Message: "There is an issue with the expiration configuration"},
		E23: ErrorContent{Code: "e23", Message: "The system has detected too many incorrect attempts"},
		E24: ErrorContent{Code: "e24", Message: "The OTP provided does not match our records"},
		E25: ErrorContent{Code: "e25", Message: "Please try to resend again later"},
		E26: ErrorContent{Code: "e26", Message: "Unable to generate recovery code"},
		E27: ErrorContent{Code: "e27", Message: "Recovery tries exceeded maximum limit"},
		E28: ErrorContent{Code: "e28", Message: "The provided recovery code is invalid"},
		E29: ErrorContent{Code: "e29", Message: "The credentials provided do not match our records"},
	}
}
