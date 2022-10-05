package user_service

type User struct {
	UUID     string `json:"uuid" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password,omitempty"`
}

type CreateUserDTO struct {
	Email          string `json:"email" validate:"email,required"`
	Password       string `json:"password" validate:"required,min=8,max=16"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password"`
	InviteId       string `json:"invite_id,omitempty"`
}
