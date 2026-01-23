package jupiter

var (
	MissingTokenProgramError       = JError{ErrorCode: "missingTokenProgram"}
	CouldNotFindAnyRouteError      = JError{ErrorCode: "couldNotFindAnyRouteError"}
	TokenNotTradableError          = JError{ErrorCode: "tokenNotTradable"}
	SlippageToleranceExceededError = JError{ErrorCode: "SlippageToleranceExceeded"}
	RequestError                   = JError{ErrorCode: "requestError"}
)

type JError struct {
	ErrorCode string
}

func (e JError) Error() string {
	return e.ErrorCode
}

func IsJError(err error) bool {

	_, ok := err.(JError)
	return ok
}
