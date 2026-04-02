package models

type User struct {
	BaseModel

	UserName     string    `json:"userName"`
	Password     string    `json:"password"`
	Email        string    `gorm:"type:varchar(100);unique_index" json:"email"`
	Phone        string    `gorm:"type:varchar(100);unique_index" json:"phone"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	UserType     string    `json:"userType" validate:"eq=PRO|eq=CLIENT"`
	ProfileImage string    `json:"profileImage"`
	Profiles     []Profile `json:"profiles" gorm:"profile"`
}
