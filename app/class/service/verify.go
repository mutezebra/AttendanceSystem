package service

import (
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"path"
	"unicode"
)

func (svc *ClassService) verifyName(v string) error {
	if len(v) <= 1 {
		return errno.New(errno.WrongClassNameLength, "class name`s len cloud not less than 2")
	}
	return nil
}

func (svc *ClassService) verifyInvitationCode(v string) error {
	if len(v) != 6 {
		return errno.New(errno.WrongInvitationCodeLength, "wrong invitation code length")
	}
	for _, i := range v {
		if !unicode.IsDigit(i) {
			return errno.New(errno.IllegalInvitationCode, "Illegal Invitation Code")
		}
	}
	return nil
}

func (svc *ClassService) verifyExcelFile(v string) error {
	if path.Ext(v) != ".xlsx" {
		return errno.New(errno.UnSupportImportFormat, "un supported import format")
	}
	return nil
}
