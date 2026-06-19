package forms

type PasswordLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"` //手机号验证自定义validator
	Password string `json:"password" form:"password" binding:"required,min=6,max=12"`
}
