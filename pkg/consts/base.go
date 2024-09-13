package consts

import base2 "github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/base/base"

var success = int32(200)
var ok = "ok"
var DefaultBase *base2.Base = &base2.Base{Code: &success, Msg: &ok}

type uid int8

var UIDKey uid

const (
	TokenKey = "CallSystem-Token"
)
