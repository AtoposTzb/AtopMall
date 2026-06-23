package forms

//这个forms是邮箱验证码的form结构体，存储邮箱
type EmailCodeForm struct {
	Email string `json:"email" form:"email" binding:"required,email"`
	Type  int    `json:"type" form:"type" binding:"required,oneof=1 2 3"`
	//Type 1 注册 2 登录 3 重置密码
}
