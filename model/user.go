package model

type CreateUserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// hide password
	Password string `json:"-"`
}

func ToUser(createReq *CreateUserReq) *User {
	return &User{Name: createReq.Name, Email: createReq.Email, Password: createReq.Password}
}

func ToCreateUserRes(user *User) *CreateUserReq {
	return &CreateUserReq{Name: user.Name, Email: user.Email, Password: user.Password}
}
