package usecase

import (
	"context"
	service "github.com/mutezebra/ClassroomRandomRollCallSystem/app/class/service"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/class"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/excel"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/pack"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/cache"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/database"
	"sync"
)

type ClassUsecase struct {
	svc   *service.ClassService
	db    *database.ClassRepository
	cache *cache.ClassCache
}

var once sync.Once
var usecase *ClassUsecase

func GetClassUsecase() *ClassUsecase {
	once.Do(func() {
		usecase = &ClassUsecase{
			svc:   service.NewClassService(),
			db:    database.NewClassRepository(),
			cache: cache.NewClassCache(),
		}
	})
	return usecase
}

func (usecase *ClassUsecase) CreateClass(ctx context.Context, req *class.CreateClassReq) (resp *class.CreateClassResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyRequest(req); err != nil {
		return nil, err
	}
	icode, classID := "", int64(0)

	if icode, classID, err = usecase.svc.CreateClass(ctx, req); err != nil {
		return nil, err
	}

	resp = new(class.CreateClassResp)
	resp.Base = consts.DefaultBase
	resp.InvitationCode, resp.ClassID = &icode, &classID
	return resp, nil
}

func (usecase *ClassUsecase) JoinClass(ctx context.Context, req *class.JoinClassReq) (resp *class.JoinClassResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyRequest(req); err != nil {
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
	if is {
		return nil, errno.New(errno.HaveInClass, "you have in class")
	}

	var in bool
	if in, err = usecase.svc.WhetherUserInClass(ctx, req.GetUID(), req.GetClassID()); err != nil {
		return nil, err
	}
	if in {
		return nil, errno.New(errno.HaveInClass, "you have in class")
	}

	var icode string
	if icode, err = usecase.svc.GetInvitationCode(ctx, req.GetClassID()); err != nil {
		return nil, err
	}
	if icode != req.GetInvitationCode() {
		return nil, errno.New(errno.WrongInvitationCode, "wrong invitation code")
	}

	if err = usecase.svc.JoinClass(ctx, req); err != nil {
		return nil, err
	}

	resp = new(class.JoinClassResp)
	resp.Base = consts.DefaultBase
	return resp, nil
}

func (usecase *ClassUsecase) ClassList(ctx context.Context, req *class.ClassListReq) (resp *class.ClassListResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyRequest(req); err != nil {
		return nil, err
	}

	var classes []*class.Class
	if classes, err = usecase.svc.ClassList(ctx, req.GetUID()); err != nil {
		return nil, err
	}

	resp = new(class.ClassListResp)
	resp.Base = consts.DefaultBase
	l := int32(len(classes))
	resp.ClassCount, resp.Classes = &l, classes
	return resp, nil
}

func (usecase *ClassUsecase) ClassStudentList(ctx context.Context, req *class.ClassStudentListReq) (resp *class.ClassStudentListResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyRequest(req); err != nil {
		return nil, err
	}

	var classID int64
	var exist bool
	if exist, classID, err = usecase.svc.WhetherUserHaveClass(ctx, req.GetUID()); err != nil {
		return nil, errno.New(errno.UserDoNotHaveClass, "user do not have class")
	}

	resp = new(class.ClassStudentListResp)
	resp.Base = consts.DefaultBase

	if !exist {
		uc := int32(0)
		resp.UserCount, resp.Students = &uc, make([]*class.StudentFormat, 0)
		return resp, nil
	}

	if req.GetClassID() == 0 {
		req.ClassID = &classID
	}
	var users []*class.StudentFormat
	if users, err = usecase.svc.ClassStudentList(ctx, req.GetClassID()); err != nil {
		return nil, err
	}

	if users, err = usecase.svc.ClassStudentListWithStatus(ctx, req.GetClassID(), users); err != nil {
		return nil, err
	}

	l := int32(len(users))
	resp.UserCount, resp.Students = &l, users
	return resp, nil
}

func (usecase *ClassUsecase) ViewInvitationCode(ctx context.Context, req *class.ViewInvitationCodeReq) (resp *class.ViewInvitationCodeResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyRequest(req); err != nil {
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
	var in bool
	if in, err = usecase.svc.WhetherUserInClass(ctx, req.GetUID(), req.GetClassID()); err != nil {
		return nil, err
	}

	if !is && !in {
		return nil, errno.New(errno.NotClassMember, "you are not class member")
	}

	var icode string
	if icode, err = usecase.svc.GetInvitationCode(ctx, req.GetClassID()); err != nil {
		return nil, err
	}

	resp = new(class.ViewInvitationCodeResp)
	resp.Base = consts.DefaultBase
	resp.InvitationCode = &icode
	return resp, nil
}

func (usecase *ClassUsecase) GetClassTeacher(ctx context.Context, req *class.GetClassTeacherReq) (resp *class.GetClassTeacherResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyRequest(req); err != nil {
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
	var in bool
	if in, err = usecase.svc.WhetherUserInClass(ctx, req.GetUID(), req.GetClassID()); err != nil {
		return nil, err
	}

	if !is && !in {
		return nil, errno.New(errno.NotClassMember, "you are not this class`s member")
	}

	var u *user.BaseUser
	if u, err = usecase.svc.GetTeacherInfo(ctx, req.GetClassID()); err != nil {
		return nil, err
	}

	resp = new(class.GetClassTeacherResp)
	resp.Base = consts.DefaultBase
	resp.Teacher = u
	return resp, nil
}

func (usecase *ClassUsecase) ImportUserAndCreateClass(ctx context.Context, req *class.ImportUserAndCreateClassReq) (resp *class.ImportUserAndCreateClassResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyRequest(req); err != nil {
		return nil, err
	}

	var users []*excel.ImportUser
	if users, err = usecase.svc.GetUserFromExcel(ctx, req.GetFile()); err != nil {
		return nil, errno.New(errno.GetUsersFromExcelFailed, "Get Users From Excel Failed")
	}

	var classID int64
	if classID, err = usecase.svc.ImportUserAndCreateClass(ctx, req.GetUID(), req.GetName(), users); err != nil {
		return nil, errno.New(errno.ImportUserFromExcelFailed, "Import User From Excel Failed")
	}

	resp = new(class.ImportUserAndCreateClassResp)
	resp.Base = consts.DefaultBase
	resp.ClassID = &classID
	return resp, nil
}

func (usecase *ClassUsecase) ChangePoint(ctx context.Context, req *class.ChangePointReq) (resp *class.ChangePointResp, err error) {
	defer func() {
		pack.LogError(err)
	}()
	if err = usecase.svc.VerifyRequest(req); err != nil {
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
		return nil, errno.New(errno.NotClassOwner, "you are not class owner")
	}

	if err = usecase.svc.ChangePoint(ctx, req.GetStuUID(), req.GetClassID(), req.GetAction(), req.GetPoint()); err != nil {
		return nil, err
	}

	resp = new(class.ChangePointResp)
	resp.Base = consts.DefaultBase
	return resp, nil
}
