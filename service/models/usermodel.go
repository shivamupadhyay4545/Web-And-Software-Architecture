package models

type Username struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
}

type Myprofile struct {
	PhotoNo   int
	Followers int
	Following int
	Photos    []Photo
}
type Details struct {
	Name string `json:"name" validate:"required,min=3,max=16"`
	Id   string `json:"id" validate:"required,min=3,max=16"`
}
