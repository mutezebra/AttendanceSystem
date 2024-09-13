package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/model/base/base"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/consts"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/errno"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/pkg/jwt"
	"net/http"
	"strings"
)

func JWT() app.HandlerFunc {
	sb := strings.Builder{}

	return func(ctx context.Context, c *app.RequestContext) {
		defer sb.Reset()
		if n, _ := sb.Write(c.GetHeader(consts.TokenKey)); n == 0 {
			code, msg := int32(errno.LackToken), "lack of token"
			c.JSON(http.StatusOK, base.BaseResp{Base: &base.Base{Code: &code, Msg: &msg}})
			c.Abort()
			return
		}

		token := sb.String()
		uid, ok, err := jwt.CheckToken(token)
		if err != nil {
			code, msg := int32(errno.WrongToken), err.Error()
			c.JSON(http.StatusOK, base.BaseResp{Base: &base.Base{Code: &code, Msg: &msg}})
			c.Abort()
			return
		}
		if !ok {
			code, msg := int32(errno.TokenExpire), "token is expired, please log in again"
			c.JSON(http.StatusOK, base.BaseResp{Base: &base.Base{Code: &code, Msg: &msg}})
			c.Abort()
			return
		}
		ctx = context.WithValue(ctx, consts.UIDKey, uid)
		c.Next(ctx)
	}
}
