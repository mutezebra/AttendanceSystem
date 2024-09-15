package usecase

import (
	"context"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/app/call/service"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/call"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/pack"
	"strconv"
	"sync"
	"time"
)

type CallUsecase struct {
	svc *service.CallService
}

var once sync.Once
var usecase *CallUsecase

func GetCallUsecase() *CallUsecase {
	once.Do(func() {
		usecase = &CallUsecase{svc: service.NewCallService()}
	})
	return usecase
}

func (usecase *CallUsecase) CallAllStudent(ctx context.Context, req *call.CallAllStudentReq) (resp *call.CallAllStudentResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyReq(req); err != nil {
		return nil, err
	}

	var exist bool
	if exist, err = usecase.svc.ClassExist(ctx, req.GetClassID()); err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.New(errno.ClassNotExist, "class not exist")
	}

	var is bool
	if is, err = usecase.svc.IsClassOwner(ctx, req.GetUID(), req.GetClassID()); err != nil {
		return nil, err
	}
	if !is {
		return nil, errno.New(errno.NotClassOwner, "not class owner")
	}

	var uids []int64
	if uids, err = usecase.svc.GetClassUids(ctx, req.GetClassID()); err != nil {
		return nil, err
	}

	var event *call.CallEvent
	if event, err = usecase.buildCallEvent(req.CallEventName, req.ClassID, req.Deadline, req.UID, uids); err != nil {
		return nil, err
	}

	var eventID int64
	if eventID, err = usecase.svc.CreateACallEvent(ctx, event); err != nil {
		return nil, err
	}

	if err = usecase.svc.CallStudents(ctx, event, uids); err != nil {
		return nil, err
	}

	resp = new(call.CallAllStudentResp)
	resp.Base = consts.DefaultBase
	resp.EventID = &eventID
	return resp, nil
}

func (usecase *CallUsecase) DoCallEvent(ctx context.Context, req *call.DoCallEventReq) (resp *call.DoCallEventResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyReq(req); err != nil {
		return nil, err
	}

	var exist bool
	if exist, err = usecase.svc.ClassExist(ctx, req.GetClassID()); err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.New(errno.ClassNotExist, "class not exist")
	}

	if exist, err = usecase.svc.EventExist(ctx, req.GetEventID()); err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.New(errno.EventNotExist, "event not exist")
	}

	if err = usecase.svc.DoCallEvent(ctx, req.GetUID(), req.GetClassID(), req.GetEventID()); err != nil {
		return nil, err
	}

	resp = new(call.DoCallEventResp)
	resp.Base = consts.DefaultBase
	return resp, nil
}

func (usecase *CallUsecase) UndoCallEvents(ctx context.Context, req *call.UndoCallEventsReq) (resp *call.UndoCallEventsResp, err error) {
	defer func() {
		pack.LogError(err)
	}()

	if err = usecase.svc.VerifyReq(req); err != nil {
		return nil, err
	}

	exist, err := usecase.svc.ClassExist(ctx, req.GetClassID())
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.New(errno.ClassNotExist, "class not exist")
	}

	resp = &call.UndoCallEventsResp{
		Base: consts.DefaultBase,
	}

	event, exist := usecase.svc.GetCallEvent(ctx, req.GetClassID())
	if !exist {
		resp.Exist = &exist
		return resp, nil
	}

	done, err := usecase.svc.WhetherUserDoEvent(ctx, req.GetUID(), req.GetClassID())
	if err != nil {
		return nil, err
	}

	if done {
		exist = false
		resp.Exist = &exist
		return resp, nil
	}

	// 填充事件信息
	id, name, classID, startTime, endTime := event.GetID(), event.GetCallEventName(), event.GetClassID(), event.GetStartTime(), event.GetEndTime()
	resp.Event = &call.CallEvent{
		ID:            &id,
		CallEventName: &name,
		ClassID:       &classID,
		StartTime:     &startTime,
		EndTime:       &endTime,
	}

	exist = true
	resp.Exist = &exist
	return resp, nil
}

func (usecase *CallUsecase) RandomCall(ctx context.Context, req *call.RandomCallReq) (resp *call.RandomCallResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyReq(req); err != nil {
		return nil, err
	}

	var exist bool
	if exist, err = usecase.svc.ClassExist(ctx, req.GetClassID()); err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.New(errno.ClassNotExist, "class not exist")
	}

	if have := usecase.svc.WhetherClassHaveEvent(req.GetClassID()); have {
		return nil, errno.New(errno.HaveCallEvent, "now class have a call event")
	}

	var is bool
	if is, err = usecase.svc.IsClassOwner(ctx, req.GetUID(), req.GetClassID()); err != nil {
		return nil, err
	}
	if !is {
		return nil, errno.New(errno.NotClassOwner, "not class owner")
	}

	var items map[int64]int
	if items, err = usecase.svc.GetUidAndWeight(ctx, req.GetClassID()); err != nil {
		return nil, err
	}

	var users []*user.BaseUser
	if users, err = usecase.svc.GetBaseUserByID(ctx, usecase.convertItems2Uids(items)); err != nil {
		return nil, err
	}

	var event *call.CallEvent
	if event, err = usecase.buildCallEvent(req.CallEventName, req.ClassID, req.Deadline, req.UID, usecase.convertItems2Uids(items)); err != nil {
		return nil, err
	}

	var eventID int64
	if eventID, err = usecase.svc.CreateACallEvent(ctx, event); err != nil {
		return nil, err
	}

	if err = usecase.svc.RandomCallUser(ctx, items, int(req.GetCallNumber()), event); err != nil {
		return nil, err
	}

	resp = new(call.RandomCallResp)
	resp.Base = consts.DefaultBase
	resp.EventID = &eventID
	for _, u := range users {
		if weighet, ok := items[u.GetUID()]; ok {
			w := int32(weighet)
			u.Weight = &w
		}
	}
	resp.Users = users

	return resp, nil
}

func (usecase *CallUsecase) HistoryCallEvent(ctx context.Context, req *call.HistoryCallEventReq) (resp *call.HistoryCallEventResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyReq(req); err != nil {
		return nil, err
	}

	var exist bool
	if exist, err = usecase.svc.ClassExist(ctx, req.GetClassID()); err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.New(errno.ClassNotExist, "class not exist")
	}

	var in bool
	if in, err = usecase.svc.WhetherUserInClass(ctx, req.GetUID(), req.GetClassID()); err != nil {
		return nil, err
	}
	if !in {
		return nil, errno.New(errno.NotClassMember, "you are not this class member")
	}

	var events []*call.CallEvent
	if events, err = usecase.svc.GetCallEvents(ctx, req.GetClassID()); err != nil {
		return nil, err
	}

	resp = new(call.HistoryCallEventResp)
	resp.Base = consts.DefaultBase
	resp.Events = events
	return resp, nil
}

func (usecase *CallUsecase) extraUidFromBaseUser(users []*user.BaseUser) []int64 {
	uids := make([]int64, 0, len(users))
	for _, u := range users {
		uids = append(uids, u.GetUID())
	}
	return uids
}

func (usecase *CallUsecase) convertUids2String(uids []int64) []string {
	strs := make([]string, 0, len(uids))
	for _, u := range uids {
		strs = append(strs, strconv.FormatInt(u, 10))
	}
	return strs
}

func (usecase *CallUsecase) buildCallEvent(name *string, classID *int64, deadline *int16, callerID *int64, uids []int64) (*call.CallEvent, error) {
	now := time.Now().Unix()
	Deadline := time.Now().Add(time.Duration(*deadline) * time.Minute).Unix()
	return &call.CallEvent{
		CallEventName: name,
		ClassID:       classID,
		CallerID:      callerID,
		StartTime:     &now,
		EndTime:       &Deadline,
	}, nil
}

func (usecase *CallUsecase) convertItems2Uids(items map[int64]int) []int64 {
	uids := make([]int64, 0, len(items))
	for k := range items {
		uids = append(uids, k)
	}
	return uids
}
