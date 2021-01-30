package common

const (
	ECONFLICT     = 409
	EINTERNAL     = 500
	EINVALID      = 400
	EUNAUTHORIZED = 401
	EFORBIDDEN    = 403
	ENOTFOUND     = 404

	EMINTERNAL = "An internal error occurred."
	EMINVALID  = "An error ocurred while processing the request."
	EMSEVERAL  = "One or more errors ocurred while processing the request."
)

type Error struct {
	Code int `json:"-"`

	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`

	Errors map[string]string `json:"errors,omitempty"`
}

func CreateGenericInternalError(err error) *Error {
	return &Error{Code: EINTERNAL, Message: EMINTERNAL, Detail: err.Error()}
}

func CreateGenericBadRequestError(err error) *Error {
	return &Error{Code: EINVALID, Message: EMINVALID, Detail: err.Error()}
}

func CreateBadRequestError(message string) *Error {
	return &Error{Code: EINVALID, Message: message}
}

func CreateConflictError(message string) *Error {
	return &Error{Code: ECONFLICT, Message: message}
}

func CreateFormError(errors map[string]string) *Error {
	return &Error{Code: EINVALID, Message: EMSEVERAL, Errors: errors}
}

func CreateNotFoundError(message string) *Error {
	return &Error{Code: ENOTFOUND, Message: message}
}
