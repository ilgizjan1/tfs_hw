package request

type CreateUserRequest struct {
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}
