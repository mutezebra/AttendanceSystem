package pack

import (
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/log"
	"github.com/pkg/errors"
)

func LogError(err error) {
	if err == nil {
		return
	}
	var e errno.Errno
	if errors.As(err, &e) {
		log.LogrusObj.Debugln(err.Error())
	} else if errors.Cause(err) != nil {
		log.LogrusObj.Errorf("stack track:\norigin error: %+v\n", err)
	} else {
		log.LogrusObj.Errorln(err.Error())
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
