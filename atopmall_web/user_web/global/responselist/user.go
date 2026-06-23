package responselist

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"nickName"`
	BirthDay string `json:"birthDay"`
	Gender   string `json:"gender"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
}
