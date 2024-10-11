package service

import (
	"context"
	"fmt"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/config"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/ems"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/utils"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/cache"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/repository/database"
	"github.com/mutezebra/tiktok/pkg/snowflake"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type UserService struct {
	db      *database.UserRepository
	classDB *database.ClassRepository
	cache   *cache.UserCache
	pnRe    *regexp.Regexp
	pwdRe   *regexp.Regexp
	ems     *ems.Client
}

func NewUserService() *UserService {
	phoneNumberRe := regexp.MustCompile("^1\\d{10}$") // phone_number
	pwdRe := regexp.MustCompile("^[a-zA-Z0-9_.]+$")   // password pattern (without checks for lower/upper case)

	return &UserService{
		db:      database.NewUserRepository(),
		classDB: database.NewClassRepository(),
		cache:   cache.NewUserCache(),
		pnRe:    phoneNumberRe,
		pwdRe:   pwdRe,
		ems:     ems.NewEmsClient(),
	}
}

func (svc *UserService) VerifyRequest(req interface{}) error {
	switch req.(type) {
	case *user.RegisterReq:
		req := req.(*user.RegisterReq)
		if err := svc.verifyName(req.GetName()); err != nil {
			return err
		}
		if err := svc.verifyPhoneNumber(req.GetPhoneNumber()); err != nil {
			return err
		}
		if err := svc.verifyStudentNumber(req.GetStudentNumber()); err != nil {
			return err
		}
		if err := svc.verifyPassword(req.GetPassword()); err != nil {
			return err
		}
		if err := svc.verifyVerifyCode(req.GetVerifyCode()); err != nil {
			return err
		}
	case *user.GetVerifyCodeReq:
		req := req.(*user.GetVerifyCodeReq)
		if err := svc.verifyPhoneNumber(req.GetPhoneNumber()); err != nil {
			return err
		}
	case *user.LoginReq:
		req := req.(*user.LoginReq)
		if err := svc.verifyPhoneNumber(req.GetPhoneNumber()); err != nil {
			return err
		}
	case *user.ChangePasswordReq:
	default:
		return errors.Wrap(errors.New("unknown req type"), "unknown req type")
	}
	return nil
}

func (svc *UserService) EncryptPassword(pwd string) (string, error) {
	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed when bcrypt %s", pwd))
	}
	return string(passwordDigest), nil
}

func (svc *UserService) CheckPassword(pwd, passwordDigest string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordDigest), []byte(pwd)) == nil
}

func (svc *UserService) PhoneNumberExist(phoneNumber string) (bool, error) {
	return svc.db.PhoneNumberExist(phoneNumber)
}

func (svc *UserService) CreateUser(req *user.RegisterReq) error {
	uid := snowflake.GenerateID(consts.UserWorkerID, consts.CenterID)
	return svc.db.CreateUser(&user.User{
		ID:             &uid,
		Name:           req.Name,
		StudentNumber:  req.StudentNumber,
		Avatar:         &config.Conf.Paths.DefaultAvatarPath,
		PhoneNumber:    req.PhoneNumber,
		PasswordDigest: req.Password,
	})
}

func (svc *UserService) SendEms(ctx context.Context, phoneNumber string) (code string, err error) {
	code = utils.GenerateCode(6)
	if err = svc.ems.SendEms(phoneNumber, code); err != nil {
		return "", err
	}
	return code, svc.cache.PutVerifyCode(ctx, phoneNumber, code)
}

func (svc *UserService) FindUserPassword(phoneNumber string) (pwd string, err error) {
	return svc.db.FindUserPassword(phoneNumber)
}

func (svc *UserService) FindUserPasswordByID(uid int64) (pwd string, err error) {
	return svc.db.FindUserPasswordByUID(uid)
}

func (svc *UserService) FindUIDByPhoneNumber(phoneNumber string) (int64, error) {
	return svc.db.FindUIDByPhoneNumber(phoneNumber)
}

func (svc *UserService) ChangePassword(uid int64, newPassword string) error {
	pwd, err := svc.EncryptPassword(newPassword)
	if err != nil {
		return err
	}
	return svc.db.ChangePassword(uid, pwd)
}

func (svc *UserService) FindUserByID(uid int64) (*user.User, error) {
	return svc.db.FindUserByID(uid)
}

func (svc *UserService) GetWeightByUID(uid int64) (int32, error) {
	classID, err := svc.classDB.WhetherUserInAClass(context.Background(), uid)
	if err != nil {
		return 0, err
	}
	if classID == 0 {
		return 0, nil
	}

	return svc.classDB.FindWeightByUIDANDClassID(context.Background(), uid, classID)
}
