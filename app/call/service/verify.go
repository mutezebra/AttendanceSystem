package service

import "github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"

func (svc *CallService) verifyCallNumber(v int64) error {
	if v <= 0 {
		return errno.New(errno.IllegalCallNumber, "call number must better than 0")
	}
	return nil
}

func (svc *CallService) verifyDeadline(v int16) error {
	if v <= 0 || v >= 60*24 {
		return errno.New(errno.IllegalDeadline, "deadline must between 1 and 60 * 24")
	}
	return nil
}

func (svc *CallService) verifyEventName(v string) error {
	if len(v) <= 1 {
		return errno.New(errno.WrongEventNameLength, "event Name`s len cloud not less than 2")
	}
	return nil
}
