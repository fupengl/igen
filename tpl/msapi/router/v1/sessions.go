package v1

import (
	"context"

	"github.com/gin-gonic/gin"

	"igen/lib/constant"
	"igen/lib/logger"
	"igen/lib/render"
	"igen/lib/util/rpcUtil"
	"igen/msdemo/middleware"
	"igen/msdemo/router/v1/helper"
	"igen/probuf/pb"
)

// CreateSession 登录
func CreateSession(c *gin.Context) {
	var form helper.SigninForm
	err := c.BindJSON(&form)
	if err != nil {
		logger.Ctx(c).Error("登录失败", logger.Err(err))
		render.Err400(c)
		return
	}

	var resp *pb.AuthModel
	var errResp error
	err = rpcUtil.AuthClient(c, func(ctx context.Context, client pb.AuthClient) {
		resp, errResp = client.Signin(ctx, &pb.AuthSigninForm{
			App:         constant.UserProject,
			Username:    form.Username,
			Password:    form.Password,
			Phone:       form.Username,
			Code:        form.ValidationCode,
			OauthUniqId: form.OpenID,
		})
	})

	err = rpcUtil.Err(c, err, errResp)
	if err != nil {
		render.Err500(c, err.Error())
		return
	}

	out := gin.H{}
	out["sessionToken"] = resp.GetToken()

	render.OK(c, out)
}

// DeleteSession 登出
func DeleteSession(c *gin.Context) {
	auth, _ := middleware.GetToken(c)

	var respErr error
	err := rpcUtil.AuthClient(c, func(ctx context.Context, client pb.AuthClient) {
		_, respErr = client.Signout(ctx, &pb.AuthToken{
			Token:    auth.Token,
			Platform: middleware.GetPlatform(c),
		})
	})

	err = rpcUtil.Err(c, err, respErr)
	if err != nil {
		logger.Ctx(c).Error("登出失败", logger.Err(err))
	}

	render.OK(c)
}
