package middleware

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"

	"igen/lib/logger"
	gMiddleware "igen/lib/middleware"
	"igen/lib/render"
	"igen/lib/util/rpcUtil"
	"igen/probuf/pb"
)

// TokenKey
const TokenKey = "gin-jwt-token"

// GetToken return session object
func GetToken(c *gin.Context) (*pb.AuthModel, error) {
	appID := GetPlatform(c)

	var am *pb.AuthModel
	var amErr error

	value, exists := c.Get(TokenKey)
	if exists {
		var ok bool
		am, ok = value.(*pb.AuthModel)
		if ok && am.Token != "" {
			return am, nil
		}
		return nil, errors.New("未登录")
	}

	token, err := gMiddleware.GetSessionToken(c)
	if err != nil {
		logger.Ctx(c).Error("未登录", logger.Err(err))
		return nil, errors.New("未登录")
	}

	err = rpcUtil.AuthClient(c, func(ctx context.Context, client pb.AuthClient) {
		am, amErr = client.Check(ctx, &pb.AuthToken{
			Token:    token,
			Platform: appID,
		})
	})

	err = rpcUtil.Err(c, err, amErr)
	if err != nil {
		logger.Ctx(c).Error("session过期", logger.Err(err))
		return nil, errors.New("session过期")
	}

	if am.Token == "" {
		logger.Ctx(c).Error("session不存在")
		return nil, errors.New("session不存在")
	}

	c.Set(TokenKey, am)
	return am, nil
}

// CheckToken check session token
func CheckToken() func(*gin.Context) {
	return func(c *gin.Context) {
		_, err := GetToken(c)
		if err != nil {
			logger.Ctx(c).Error("session错误", logger.Err(err))
			render.Abort(c, 401, "session错误，请重新登录")
			return
		}

		c.Next()
	}
}
