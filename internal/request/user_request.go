package request

type UserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
