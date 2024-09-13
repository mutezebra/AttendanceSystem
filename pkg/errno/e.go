package errno

type e int32

// Basic
var (
	LackToken   e = 1000
	TokenExpire e = 1001
	WrongToken  e = 1001
)

// User
var (
	EmptyUserName        e = 10000
	EmptyPassword        e = 10001
	EmptyPhoneNumber     e = 10002
	EmptyStudentNumber   e = 10003
	UserNameOutOfLen     e = 10004
	SpaceUserName        e = 10005
	IllegalPhoneNumber   e = 10006
	IllegalStudentNumber e = 10007
	IllegalPasswordLen   e = 10008
	IllegalPassword      e = 10009
	ExistPhoneNumber     e = 10010
	VerifyCodeLen        e = 10011
	IllegalVerifyCode    e = 10012
	VerifyCodeExpired    e = 10012
	WrongVerifyCode      e = 10013
	WrongPassword        e = 10014

	UserRegisterFailed e = 10100
)
