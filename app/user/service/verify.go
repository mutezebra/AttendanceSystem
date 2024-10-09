package service

import (
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"unicode"
)

func (svc *UserService) verifyName(v string) errno.Errno {
	name := v
	if len(name) == 0 {
		return errno.New(errno.EmptyUserName, "user_name should not be empty")
	}
	if len(name) >= 20 {
		return errno.New(errno.UserNameOutOfLen, "user_name should be limited to 20 characters")
	}
	for i := range name {
		if name[i] == ' ' || name[i] == '\n' || name[i] == '\t' {
			return errno.New(errno.SpaceUserName, "there should be no white space or line breaks in the user name")
		}
	}
	return nil
}

func (svc *UserService) verifyPhoneNumber(v string) errno.Errno {
	phone := v
	if len(phone) == 0 {
		return errno.New(errno.EmptyPhoneNumber, "phone_number should not be empty")
	}
	if !svc.pnRe.MatchString(phone) {
		return errno.New(errno.IllegalPhoneNumber, "illegal phone_number")
	}
	return nil
}

func (svc *UserService) verifyStudentNumber(v string) errno.Errno {
	number := v
	if len(number) == 0 {
		return errno.New(errno.EmptyStudentNumber, "student_number should not be empty")
	}
	ok := true
	for _, c := range number {
		if !unicode.IsDigit(c) {
			ok = false
		}
	}
	if !ok {
		return errno.New(errno.IllegalStudentNumber, "student_number should consists with numbers")
	}
	return nil
}

func (svc *UserService) verifyPassword(v string) errno.Errno {
	pwd := v
	if len(pwd) > 15 || len(pwd) < 6 {
		return errno.New(errno.IllegalPasswordLen, "password`s len should be between 6 and 20")
	}

	if !svc.pwdRe.MatchString(pwd) {
		return errno.New(errno.IllegalPassword, "password can only contain letters, numbers, underscores, and periods")
	}

	return nil
}

func (svc *UserService) verifyVerifyCode(v string) errno.Errno {
	verifyCode := v
	if len(verifyCode) != 6 {
		return errno.New(errno.VerifyCodeLen, "the code length should be 6")
	}
	for _, i := range verifyCode {
		if !unicode.IsDigit(i) {
			return errno.New(errno.IllegalVerifyCode, "verify should have 6 number")
		}
	}
	return nil
}
