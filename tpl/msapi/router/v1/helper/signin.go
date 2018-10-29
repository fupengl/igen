package helper

//SigninForm 登录表单
type SigninForm struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	ValidationCode string `json:"validationCode"`
	OpenID         string `json:"openId"`
}

type SigninResp struct {
	ID   string      `json:"id"`
	User interface{} `json:"user"`
}
