package v1

import (
	"context"

	"github.com/gin-gonic/gin"

	"igen/lib/constant"
	"igen/lib/logger"
	"igen/lib/render"
	"igen/lib/util/rpcUtil"
	"igen/probuf/pb"
)

// CreateSMS 发送手机验证码
func CreateSMS(c *gin.Context) {
	var form struct {
		Type  int64
		Phone string
	}

	err := c.BindJSON(&form)
	if err != nil {
		logger.Ctx(c).Error(err.Error())
		render.Err400(c)
		return
	}

	var resp *pb.Bool
	var errResp error
	err = rpcUtil.AuthClient(c, func(ctx context.Context, client pb.AuthClient) {
		resp, errResp = client.NewCaptcha(ctx, &pb.AuthCaptchaForm{
			App:   constant.UserProject,
			Type:  form.Type,
			Phone: form.Phone,
		})
	})

	err = rpcUtil.Err(c, err, errResp)
	if err != nil {
		render.Err500(c, "服务器错误")
		return
	}

	render.OK(c)
}
