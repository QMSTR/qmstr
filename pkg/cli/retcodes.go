package cli

const (
	ReturnCodeSuccess = iota
	ReturnCodeCliError
	ReturnCodeServerCommunicationError
	ReturnCodeResponseFalseError
	ReturnCodeTimeout
	ReturnCodeParameterError
	ReturnCodeSysError
	ReturnCodeFormatError
)
