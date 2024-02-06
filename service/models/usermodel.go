package models

type Username struct {
	Username string `json:"name" validate:"required,min=3,max=16"`
}

type Myprofile struct {
	PhotoNo   int
	Followers int
	Following int
	Photos    []Photo
}
