package response

type UserLoginResponse struct {
	Token    string `json:"token"`
	UserName string `json:"user_name"`
	UserId   int64  `json:"user_id"`
}
