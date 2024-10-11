package service

import (
	"context"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/call"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/pack"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/weightedrand"
	"time"
)

func (svc *CallService) CallStudents(ctx context.Context, event *call.CallEvent, uids []int64) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	if exist := svc.whetherClassHaveEvent(event.GetClassID()); exist {
		return errno.New(errno.HaveCallEvent, "now class have a call event")
	}
	if err := svc.cache.SetNewEvent(ctx, event.GetClassID(), uids); err != nil {
		return err
	}
	if err := svc.setCallEventDeadline(ctx, event); err != nil {
		return err
	}
	return nil
}

func (svc *CallService) RandomCallUser(ctx context.Context, items map[int64]int, callCount int, action, number int8) ([]int64, error) {
	itemList := make([]*weightedrand.Item, 0, len(items))
	for k, v := range items {
		itemList = append(itemList, &weightedrand.Item{Key: k, Weight: v})
	}
	return weightedrand.WeightedRandom(itemList, callCount, action, number)
}

func (svc *CallService) DoCallEvent(ctx context.Context, uid, classID, eventID int64) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	var exist bool
	var e *SvcCallEvent
	if e, exist = svc.getCallEvent(ctx, classID); !exist {
		return errno.New(errno.ExpireORNotExist, "event expire or not exist")
	}
	if e.EventID != eventID {
		return errno.New(errno.ExpiredEvent, "event_id not match")
	}

	var err error
	if exist, err = svc.cache.WhetherUserExist(ctx, classID, uid); err != nil {
		return err
	}
	if !exist {
		return errno.New(errno.NotClassMemberORNewMember, "not class member or join in class after call event")
	}

	if err = svc.cache.AddDoneUser(ctx, classID, []int64{uid}); err != nil {
		return err
	}
	return nil
}

func (svc *CallService) GetUndoUids(ctx context.Context, classID int64) ([]int64, error) {
	return svc.cache.GetUidDiff(ctx, classID)
}

func (svc *CallService) GetDoneUids(ctx context.Context, classID int64) ([]int64, error) {
	return svc.cache.GetUidInter(ctx, classID)
}

func (svc *CallService) UpdateEventDone(ctx context.Context, eventID, classID int64, done, undo []int64) func() {
	return func() {
		pack.LogError(svc.db.UpdateEventDone(ctx, eventID, classID, done, undo))
	}
}

func (svc *CallService) WhetherClassHaveEvent(classID int64) bool {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	return svc.whetherClassHaveEvent(classID)
}

func (svc *CallService) GetCallEvent(ctx context.Context, classID int64) (*SvcCallEvent, bool) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	return svc.getCallEvent(ctx, classID)
}

func (svc *CallService) whetherClassHaveEvent(classID int64) bool {
	_, have := svc.unDoneEvents[classID]
	return have
}

func (svc *CallService) getCallEvent(ctx context.Context, classID int64) (*SvcCallEvent, bool) {
	if e, have := svc.unDoneEvents[classID]; have {
		return e, have
	}
	return nil, false
}

func (svc *CallService) setCallEventDeadline(ctx context.Context, event *call.CallEvent) error {
	callE := &SvcCallEvent{
		EventID:   event.GetID(),
		Name:      event.GetCallEventName(),
		ClassID:   event.GetClassID(),
		StartTime: event.GetStartTime(),
		EndTime:   event.GetEndTime(),
	}

	svc.unDoneEvents[callE.ClassID] = callE
	if err := svc.cache.SaveSvcCallEvent(ctx, callE.ClassID, callE.convertToJson(), time.Duration(callE.EndTime-callE.StartTime)*time.Second); err != nil {
		return err
	}
	return nil
}

func (svc *CallService) getTheTimedEvents() []*SvcCallEvent {
	now := time.Now().Unix()
	events := make([]*SvcCallEvent, 0)
	for _, v := range svc.unDoneEvents {
		if now > v.EndTime {
			events = append(events, v)
		}
	}
	for _, v := range events {
		delete(svc.unDoneEvents, v.ClassID)
	}
	return events
}
