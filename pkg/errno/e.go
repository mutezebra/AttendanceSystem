package errno

type e int32

// Basic
var (
	LackToken   e = 1000
	TokenExpire e = 1001
	WrongToken  e = 1002
)

// User
const (
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
	UnExistPhoneNumber   e = 10015
	ChangePasswordFailed e = 10016
)

// class
const (
	WrongClassNameLength      e = 20000
	ClassNotExist             e = 20001
	WrongInvitationCodeLength e = 20002
	WrongInvitationCode       e = 20003
	IllegalInvitationCode     e = 20004
	NotClassMember            e = 20005
	HaveInClass               e = 20006
	NotClassOwner             e = 20007
	WrongExcelFormat          e = 20008
	UnSupportImportFormat     e = 20009
	GetUsersFromExcelFailed   e = 20010
	ImportUserFromExcelFailed e = 20011
	UserDoNotHaveClass        e = 20012
)

// call
const (
	IllegalCallNumber         e = 30000
	IllegalDeadline           e = 30001
	HaveCallEvent             e = 30002
	WrongEventNameLength      e = 30003
	EventNotExist             e = 30004
	ExpireORNotExist          e = 30005
	ExpiredEvent              e = 30006
	NotClassMemberORNewMember e = 30007
	NotEnoughStudent          e = 30008
)
