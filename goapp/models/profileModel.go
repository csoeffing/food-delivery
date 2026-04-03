package models

type Profile struct {
	BaseModel

	UserID       int          `json:"userId"`
	FirstName    string       `json:"firstName"`
	LastName     string       `json:"lastName"`
	ProfileImage string       `json:"profileImage"`
	UserType     string       `json:"userType" validate:"eq=PRO|eq=CLIENT"`
	ProType      string       `json:"proType"  validate:"eq=CHEF|eq=RIDER"`
	UserName     string       `json:"userName"`
	Restaurant   []Restaurant `json:"restaurant"`
}
