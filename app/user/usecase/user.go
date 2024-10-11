package usecase

import (
	"context"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/app/user/service"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/jwt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/pack"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/cache"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/database"
	"sync"
)

type UserUsecase struct {
	svc   *service.UserService
	db    *database.UserRepository
	cache *cache.UserCache
}

var once sync.Once
var usecase *UserUsecase

func GetUserUsecase() *UserUsecase {
	once.Do(func() {
		usecase = &UserUsecase{
			svc:   service.NewUserService(),
			db:    database.NewUserRepository(),
			cache: cache.NewUserCache(),
		}
	})
	return usecase
}

func (usecase *UserUsecase) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	defer func() {
		pack.LogError(err)
	}()

	if err = usecase.svc.VerifyRequest(req); err != nil {
		return nil, err
	}

	exist, code := false, ""
	if exist, code, err = usecase.cache.WhetherVerifyCodeExist(ctx, req.GetPhoneNumber()); err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.New(errno.VerifyCodeExpired, "please retry to send the verify code")
	}
	if code != req.GetVerifyCode() {
		return nil, errno.New(errno.WrongPassword, "wrong verify code")
	}

	if exist, err = usecase.svc.PhoneNumberExist(req.GetPhoneNumber()); err != nil {
		return nil, err
	}
	if exist {
		return nil, errno.New(errno.ExistPhoneNumber, "phone number have exist")
	}

	var passwordDigest string
	if passwordDigest, err = usecase.svc.EncryptPassword(req.GetPassword()); err != nil {
		return nil, err
	}
	req.Password = &passwordDigest

	if err = usecase.svc.CreateUser(req); err != nil {
		return nil, err
	}

	resp = new(user.RegisterResp)
	resp.Base = consts.DefaultBase
	return resp, nil
}

func (usecase *UserUsecase) GetVerifyCode(ctx context.Context, req *user.GetVerifyCodeReq) (resp *user.GetVerifyCodeResp, err error) {
	defer func() {
		pack.LogError(err)
	}()

	if err = usecase.svc.VerifyRequest(req); err != nil {
		return nil, err
	}

	var exist bool
	if exist, err = usecase.svc.PhoneNumberExist(req.GetPhoneNumber()); err != nil {
		return nil, err
	}
	if exist {
		return nil, errno.New(errno.ExistPhoneNumber, "phone number have exist")
	}

	var code string
	if code, err = usecase.svc.SendEms(ctx, req.GetPhoneNumber()); err != nil {
		return nil, err
	}

	resp = new(user.GetVerifyCodeResp)
	resp.Base = consts.DefaultBase
	resp.VerifyCode = &code
	return resp, nil
}

func (usecase *UserUsecase) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	defer func() {
		pack.LogError(err)
	}()

	if err = usecase.svc.VerifyRequest(req); err != nil {
		return nil, err
	}
	var uid int64
	if uid, err = usecase.svc.FindUIDByPhoneNumber(req.GetPhoneNumber()); err != nil {
		return nil, errno.New(errno.UnExistPhoneNumber, "手机号不存在")
	}

	var pwdDigest string
	if pwdDigest, err = usecase.svc.FindUserPassword(req.GetPhoneNumber()); err != nil {
		return nil, err
	}

	if ok := usecase.svc.CheckPassword(req.GetPassword(), pwdDigest); !ok {
		return nil, errno.New(errno.WrongPassword, "wrong password")
	}

	var token string
	if token, err = jwt.GenerateToken(uid); err != nil {
		return nil, err
	}

	resp = new(user.LoginResp)
	resp.Base = consts.DefaultBase
	resp.Token = &token
	return resp, nil
}

func (usecase *UserUsecase) ChangePassword(ctx context.Context, req *user.ChangePasswordReq) (resp *user.ChangePasswordResp, err error) {
	defer func() {
		pack.LogError(err)
	}()

	if err = usecase.svc.VerifyRequest(req); err != nil {
		return nil, err
	}

	var pwdDigest string
	if pwdDigest, err = usecase.svc.FindUserPasswordByID(req.GetUID()); err != nil {
		return nil, err
	}

	if ok := usecase.svc.CheckPassword(req.GetOldPassword(), pwdDigest); !ok {
		return nil, errno.New(errno.WrongPassword, "wrong password")
	}

	if err = usecase.svc.ChangePassword(req.GetUID(), req.GetNewPassword()); err != nil {
		return nil, errno.New(errno.ChangePasswordFailed, "change password failed")
	}
	resp = new(user.ChangePasswordResp)
	resp.Base = consts.DefaultBase
	return resp, nil
}

func (usecase *UserUsecase) UserInfo(ctx context.Context, req *user.UserInfoReq) (resp *user.UserInfoResp, err error) {
	defer func() {
		pack.LogError(err)
	}()

	var u *user.User
	if u, err = usecase.svc.FindUserByID(req.GetUID()); err != nil {
		return nil, err
	}

	var weight int32
	if weight, err = usecase.svc.GetWeightByUID(req.GetUID()); err != nil {
		return nil, err
	}

	resp = new(user.UserInfoResp)
	resp.Base = consts.DefaultBase
	resp.UID = u.ID
	resp.Name = u.Name
	resp.StudentNumber = u.StudentNumber
	resp.PhoneNumber = u.PhoneNumber
	resp.Point = &weight
	return resp, nil
}
