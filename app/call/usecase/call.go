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

	var classID int64
	if _, classID, err = usecase.svc.WhetherUserHaveClass(ctx, req.GetUID()); err != nil {
		return nil, errno.New(errno.UserDoNotHaveClass, "user do not have class")
	}
	if req.GetClassID() == 0 {
		req.ClassID = &classID
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

	if err = usecase.svc.VerifyReq(req); err != nil {
		return nil, err
	}

	var classID int64
	var exist bool
	if exist, classID, err = usecase.svc.WhetherUserInAClass(ctx, req.GetUID()); err != nil {
		return nil, errno.New(errno.UserDoNotHaveClass, "get user class error")
	}
	if !exist {
		return nil, errno.New(errno.UserDoNotHaveClass, "user do not in a class")
	}
	if req.GetClassID() == 0 {
		req.ClassID = &classID
	}

	var eventID int64
	if eventID = usecase.svc.GetEventIDFromClass(req.GetClassID()); eventID == 0 {
		return nil, errno.New(errno.EventNotExist, "do not have event")
	}
	req.EventID = &eventID

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

	var classID int64
	var exist bool
	if exist, classID, err = usecase.svc.WhetherUserInAClass(ctx, req.GetUID()); err != nil {
		return nil, errno.New(errno.UserDoNotHaveClass, "get user class error")
	}
	if !exist {
		return nil, errno.New(errno.UserDoNotHaveClass, "user do not in a class")
	}
	if req.GetClassID() == 0 {
		req.ClassID = &classID
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

	var classID int64
	var exist bool
	if exist, classID, err = usecase.svc.WhetherUserHaveClass(ctx, req.GetUID()); err != nil {
		return nil, errno.New(errno.UserDoNotHaveClass, "user do not have class")
	}
	if !exist {
		return nil, errno.New(errno.UserDoNotHaveClass, "user do not have class")
	}

	if req.GetClassID() == 0 {
		req.ClassID = &classID
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

	var uids []int64
	if uids, err = usecase.svc.RandomCallUser(ctx, items, int(req.GetCallNumber()), req.GetAction(), req.GetNumber()); err != nil {
		return nil, err
	}

	result := make([]*user.BaseUser, 0, len(uids))
	for i, u := range users {
		for _, id := range uids {
			if id == u.GetUID() {
				result = append(result, users[i])
			}
		}
	}

	resp = new(call.RandomCallResp)
	resp.Base = consts.DefaultBase
	for _, u := range result {
		if weighet, ok := items[u.GetUID()]; ok {
			w := int32(weighet)
			u.Weight = &w
		}
	}
	resp.Users = result

	return resp, nil
}

func (usecase *CallUsecase) HistoryCallEvent(ctx context.Context, req *call.HistoryCallEventReq) (resp *call.HistoryCallEventResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyReq(req); err != nil {
		return nil, err
	}

	var classID int64
	var exist bool
	if exist, classID, err = usecase.svc.WhetherUserInAClass(ctx, req.GetUID()); err != nil {
		return nil, errno.New(errno.UserDoNotHaveClass, "get user class error")
	}
	if !exist {
		return nil, errno.New(errno.UserDoNotHaveClass, "user do not in a class")
	}
	if req.GetClassID() == 0 {
		req.ClassID = &classID
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
	if *deadline == -1 {
		Deadline = time.Now().Add(1 * time.Second).Unix()
	}
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
