package forms

//这个forms是密码登录的form结构体，存储手机号和密码和验证码
type PasswordLoginForm struct {
	Mobile    string `json:"mobile" form:"mobile" binding:"required,mobile"` //手机号验证自定义validator
	Password  string `json:"password" form:"password" binding:"required,min=6,max=12"`
	Captcha   string `json:"captcha" form:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `json:"captcha_id" form:"captcha_id" binding:"required"`
}

//注册的form结构体，存储手机号和密码和验证码和邮箱
type RegisterForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"` //手机号验证自定义validator
	Password string `json:"password" form:"password" binding:"required,min=6,max=12"`
	Email    string `json:"email" form:"email" binding:"required,email"`     //注册获取验证码用的邮箱
	Code     string `json:"code" form:"code" binding:"required,min=6,max=6"` //注册获取验证码用的验证码
}
