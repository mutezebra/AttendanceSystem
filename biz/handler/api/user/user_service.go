// Code generated by hertz generator.

package user

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/app/user/usecase"
	user "github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/api/user"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/base/base"
	consts2 "github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/pack"
)

// Register .
// @router /register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.RegisterReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := usecase.GetUserUsecase().Register(ctx, &req)
	if err != nil {
		resp = new(user.RegisterResp)
		httpcode, errno := pack.ProcessError(err)
		code, msg := errno.Code(), errno.Error()
		resp.Base = &base.Base{Code: &code, Msg: &msg}
		c.JSON(httpcode, resp)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// GetVerifyCode .
// @router /get-verifycode [GET]
func GetVerifyCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.GetVerifyCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := usecase.GetUserUsecase().GetVerifyCode(ctx, &req)
	if err != nil {
		resp = new(user.GetVerifyCodeResp)
		httpcode, errno := pack.ProcessError(err)
		code, msg := errno.Code(), errno.Error()
		resp.Base = &base.Base{Code: &code, Msg: &msg}
		c.JSON(httpcode, resp)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// Login .
// @router /auth/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.LoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp, err := usecase.GetUserUsecase().Login(ctx, &req)
	if err != nil {
		resp = new(user.LoginResp)
		httpcode, errno := pack.ProcessError(err)
		code, msg := errno.Code(), errno.Error()
		resp.Base = &base.Base{Code: &code, Msg: &msg}
		c.JSON(httpcode, resp)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// ChangePassword .
// @router /user/auth/change-password [POST]
func ChangePassword(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.ChangePasswordReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	uid := ctx.Value(consts2.UIDKey).(int64)
	req.UID = &uid

	resp, err := usecase.GetUserUsecase().ChangePassword(ctx, &req)
	if err != nil {
		resp = new(user.ChangePasswordResp)
		httpcode, errno := pack.ProcessError(err)
		code, msg := errno.Code(), errno.Error()
		resp.Base = &base.Base{Code: &code, Msg: &msg}
		c.JSON(httpcode, resp)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// UserInfo .
// @router /user/info [POST]
func UserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserInfoReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	uid := ctx.Value(consts2.UIDKey).(int64)
	req.UID = &uid

	resp, err := usecase.GetUserUsecase().UserInfo(ctx, &req)
	if err != nil {
		resp = new(user.UserInfoResp)
		httpcode, errno := pack.ProcessError(err)
		code, msg := errno.Code(), errno.Error()
		resp.Base = &base.Base{Code: &code, Msg: &msg}
		c.JSON(httpcode, resp)
		return
	}

	c.JSON(consts.StatusOK, resp)
}
