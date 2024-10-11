package service

import (
	"context"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/call"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/class"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/cache"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/database"
	"github.com/mutezebra/tiktok/pkg/snowflake"
	"github.com/pkg/errors"
	"sync"
)

type CallService struct {
	classDB      *database.ClassRepository
	userDB       *database.UserRepository
	db           *database.CallRepository
	cache        *cache.CallCache
	mu           sync.Mutex
	unDoneEvents map[int64]*SvcCallEvent // class -> timeStamp
}

func NewCallService() *CallService {
	svc := &CallService{
		classDB:      database.NewClassRepository(),
		userDB:       database.NewUserRepository(),
		db:           database.NewCallRepository(),
		cache:        cache.NewCallCache(),
		mu:           sync.Mutex{},
		unDoneEvents: make(map[int64]*SvcCallEvent),
	}
	svc.InitCallSystem()
	return svc
}

func (svc *CallService) VerifyReq(req interface{}) error {
	switch req.(type) {
	case *call.CallAllStudentReq:
		req := req.(*call.CallAllStudentReq)
		if err := svc.verifyDeadline(req.GetDeadline()); err != nil {
			return err
		}

	case *call.DoCallEventReq:
	case *call.UndoCallEventsReq:
	case *call.RandomCallReq:
		req := (req).(*call.RandomCallReq)
		if err := svc.verifyCallNumber(req.GetCallNumber()); err != nil {
			return err
		}
		//if err := svc.verifyDeadline(req.GetDeadline()); err != nil {
		//	return err
		//}
	case *call.HistoryCallEventReq:
	default:
		return errors.Wrap(fmt.Errorf("unknown request type"), "")
	}
	return nil
}

func (svc *CallService) ClassExist(ctx context.Context, classID int64) (bool, error) {
	return svc.classDB.ClassExistByID(ctx, classID)
}

func (svc *CallService) EventExist(ctx context.Context, eventID int64) (bool, error) {
	return svc.db.EventExist(ctx, eventID)
}

func (svc *CallService) GetAllClassStudents(ctx context.Context, classID int64) ([]*class.StudentFormat, error) {
	return svc.classDB.ClassStudentList(ctx, classID)
}

func (svc *CallService) CreateACallEvent(ctx context.Context, event *call.CallEvent) (int64, error) {
	id := snowflake.GenerateID(consts.CallWorkerID, consts.CenterID)
	event.ID = &id

	return id, svc.db.CreateCallEvent(ctx, event)
}

func (svc *CallService) IsClassOwner(ctx context.Context, uid, classID int64) (bool, error) {
	return svc.classDB.IsClassOwner(ctx, uid, classID)
}

func (svc *CallService) WhetherUserDoEvent(ctx context.Context, uid, classID int64) (bool, error) {
	var exist bool
	var err error
	if exist, err = svc.cache.WhetherUserExist(ctx, classID, uid); err != nil {
		return true, err
	}
	if !exist {
		return true, nil
	}

	var done bool
	if done, err = svc.cache.WhetherUserHaveDone(ctx, classID, uid); err != nil {
		return true, err
	}

	return done, nil
}

func (svc *CallService) GetUidAndWeight(ctx context.Context, classID int64) (map[int64]int, error) {
	return svc.db.GetUidAndWeight(ctx, classID)
}

func (svc *CallService) GetClassUids(ctx context.Context, classID int64) ([]int64, error) {
	return svc.db.GetClassUids(ctx, classID)
}

func (svc *CallService) GetCallEvents(ctx context.Context, classID int64) ([]*call.CallEvent, error) {
	return svc.db.GetCallEvents(ctx, classID)
}

func (svc *CallService) WhetherUserInClass(ctx context.Context, uid, classID int64) (bool, error) {
	return svc.classDB.WhetherUserInClass(ctx, uid, classID)
}

func (svc *CallService) GetBaseUserByID(ctx context.Context, uids []int64) ([]*user.BaseUser, error) {
	return svc.userDB.FindUserByUID(uids)
}

func (svc *CallService) WhetherUserHaveClass(ctx context.Context, uid int64) (bool, int64, error) {
	classID, err := svc.classDB.WhetherUserHaveClass(ctx, uid)
	return classID != 0, classID, err
}

func (svc *CallService) WhetherUserInAClass(ctx context.Context, uid int64) (bool, int64, error) {
	classID, err := svc.classDB.WhetherUserInAClass(ctx, uid)

	return classID != 0, classID, err
}

func (svc *CallService) GetEventIDFromClass(classID int64) int64 {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	if event, ok := svc.unDoneEvents[classID]; ok {
		return event.GetID()
	}
	return 0
}
