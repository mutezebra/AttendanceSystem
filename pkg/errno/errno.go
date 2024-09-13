package errno

type Errno interface {
	Code() int32
	Error() string
}

type errno struct {
	msg  string
	code int32
}

var Success = New(e(200), "ok")
var Unknown = New(e(500), "unknown error")

func (e *errno) Code() int32 {
	return e.code
}

func (e *errno) Error() string {
	return e.msg
}

func New(code e, msg string) Errno {
	return &errno{code: int32(code), msg: msg}
}
