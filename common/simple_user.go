package common

/* 
! For owner field in item, not showing private user info
*/
type SimpleUser struct {
	SQLModel 
	LastName string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string  `json:"first_name" gorm:"column:first_name;"`
	Status 	 int 			`json:"status" gorm:"column:status;"`
}


func (SimpleUser) TableName() string {
	return "users"
}

func (u *SimpleUser) Mask() {
	u.SQLModel.Mask(DbTypeUser)
}