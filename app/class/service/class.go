package service

import (
	"context"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/class"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/excel"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/utils"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/cache"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/database"
	"github.com/mutezebra/tiktok/pkg/snowflake"
	"github.com/pkg/errors"
)

type ClassService struct {
	db    *database.ClassRepository
	cache *cache.ClassCache
}

func NewClassService() *ClassService {
	return &ClassService{
		db:    database.NewClassRepository(),
		cache: cache.NewClassCache(),
	}
}

func (svc *ClassService) VerifyRequest(req interface{}) error {
	switch req.(type) {
	case *class.CreateClassReq:
		req := req.(*class.CreateClassReq)
		if err := svc.verifyName(req.GetName()); err != nil {
			return err
		}
	case *class.JoinClassReq:
		req := req.(*class.JoinClassReq)
		if err := svc.verifyInvitationCode(req.GetInvitationCode()); err != nil {
			return err
		}
	case *class.ClassListReq:
	case *class.ClassStudentListReq:
	case *class.ViewInvitationCodeReq:
	case *class.GetClassTeacherReq:
	case *class.ImportUserAndCreateClassReq:
		req := req.(*class.ImportUserAndCreateClassReq)
		if err := svc.verifyExcelFile(req.GetFileName()); err != nil {
			return err
		}
	case *class.ChangePointReq:

	default:
		return errors.Wrap(errors.New("unknown req type"), "unknown req type")
	}
	return nil
}

func (svc *ClassService) CreateClass(ctx context.Context, req *class.CreateClassReq) (string, int64, error) {
	id, name, count, code := snowflake.GenerateID(consts.ClassWorkerID, consts.CenterID), req.GetName(), int32(1), utils.GenerateCode(6)
	c := &class.Class{
		ID:             &id,
		Name:           &name,
		UserCount:      &count,
		InvitationCode: &code,
	}
	return code, id, svc.db.CreateClass(ctx, c, req.GetUID())
}

func (svc *ClassService) ClassExist(ctx context.Context, classID int64) (bool, error) {
	return svc.db.ClassExistByID(ctx, classID)
}

func (svc *ClassService) GetInvitationCode(ctx context.Context, classID int64) (string, error) {
	return svc.db.FindClassInvitationCode(ctx, classID)
}

func (svc *ClassService) JoinClass(ctx context.Context, req *class.JoinClassReq) error {
	return svc.db.JoinClass(ctx, req.GetUID(), req.GetClassID())
}

func (svc *ClassService) ClassList(ctx context.Context, uid int64) ([]*class.Class, error) {
	return svc.db.ClassList(ctx, uid)
}

func (svc *ClassService) IsClassOwner(ctx context.Context, uid, classID int64) (bool, error) {
	return svc.db.IsClassOwner(ctx, uid, classID)
}

func (svc *ClassService) WhetherUserInClass(ctx context.Context, uid, classID int64) (bool, error) {
	return svc.db.WhetherUserInClass(ctx, uid, classID)
}

func (svc *ClassService) ClassStudentList(ctx context.Context, classID int64) ([]*class.StudentFormat, error) {
	return svc.db.ClassStudentList(ctx, classID)
}

func (svc *ClassService) ClassStudentListWithStatus(ctx context.Context, classID int64, users []*class.StudentFormat) ([]*class.StudentFormat, error) {
	eventID, err := svc.db.RecentEvent(ctx, classID)
	if err != nil {
		return nil, err
	}

	if eventID == 0 {
		var emptyStr = ""
		for _, u := range users {
			u.Status = &emptyStr
		}
		return users, nil
	}

	var exist bool
	if exist, err = svc.cache.WhetherEventExist(ctx, classID); err != nil {
		return nil, err
	}
	var done, undone = "已签到", "未签到"
	if exist {
		var doneUsers []int64
		if doneUsers, err = svc.cache.GetUidInter(ctx, classID); err != nil {
			return nil, err
		}

		for _, u := range users {
			if svc.contain(u.GetUID(), doneUsers) {
				u.Status = &done
			} else {
				u.Status = &undone
			}
		}
		return users, nil
	}

	var undoneUsers []int64
	if undoneUsers, err = svc.db.GetUnDoneUiDS(ctx, eventID); err != nil {
		return nil, err
	}

	for _, u := range users {
		if svc.contain(u.GetUID(), undoneUsers) {
			u.Status = &undone
		} else {
			u.Status = &done
		}
	}
	return users, nil
}

func (svc *ClassService) contain(uid int64, users []int64) bool {
	for _, u := range users {
		if u == uid {
			return true
		}
	}
	return false
}

func (svc *ClassService) GetTeacherInfo(ctx context.Context, classID int64) (*user.BaseUser, error) {
	return svc.db.GetTeacherInfo(ctx, classID)
}

func (svc *ClassService) GetUserFromExcel(ctx context.Context, data []byte) ([]*excel.ImportUser, error) {
	return excel.ReadExcelToUsers(data)
}

func (svc *ClassService) ImportUserAndCreateClass(ctx context.Context, uid int64, className string, users []*excel.ImportUser) (int64, error) {
	classID := svc.cache.GetClassID(context.Background())
	icode := utils.GenerateCode(6)
	pwd := "$2a$10$dh0XpDhNdm8yFEPjhiGahukXN.BLvSs.W39AhivOPw7H3CTkbyT12"
	for _, u := range users {
		u.UID = svc.cache.GetUID(context.Background())
	}

	return classID, svc.db.ImportUserAndCreateClass(ctx, classID, uid, className, icode, pwd, users)
}

func (svc *ClassService) WhetherUserHaveClass(ctx context.Context, uid int64) (bool, int64, error) {
	classID, err := svc.db.WhetherUserHaveClass(ctx, uid)
	return classID != 0, classID, err
}

func (svc *ClassService) ChangePoint(ctx context.Context, uid, classID int64, action int8, point int32) error {
	if action == 0 {
		point = -point
	}
	return svc.db.ChangePoint(ctx, uid, classID, point)
}
