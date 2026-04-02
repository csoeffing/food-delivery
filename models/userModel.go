package models

type User struct {
	BaseModel

	UserName     string    `form:"userName" json:"userName"`
	Password     string    `form:"password" json:"password"`
	Email        string    `form:"email" gorm:"type:varchar(100);unique_index" json:"email"`
	Phone        string    `form:"phone" gorm:"type:varchar(100);unique_index" json:"phone"`
	FirstName    string    `form:"firstName" json:"firstName"`
	LastName     string    `form:"lastName" json:"lastName"`
	UserType     string    `form:"userType" json:"userType" validate:"eq=PRO|eq=CLIENT"`
	ProfileImage string    `form:"-" json:"profileImage"`
	Profiles     []Profile `form:"-" json:"profiles" gorm:"profile"`
}
