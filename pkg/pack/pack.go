package pack

import (
	"errors"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/log"
)

func LogError(err error) {
	if err == nil {
		return
	}
	var e errno.Errno
	if errors.As(err, &e) {
		log.LogrusObj.Debugln(err.Error())
	} else {
		log.LogrusObj.Errorln(errors.Unwrap(err))
	}
}

func ProcessError(err error) errno.Errno {
	if err == nil {
		return errno.Success
	}
	var e errno.Errno
	if !errors.As(err, &e) {
		return errno.Unknown
	}
	return e
}
