package model

import (
    "database/sql/driver"
    "errors"
    "fmt"
    "first-proj/common"
)

const EntityName = "User"

type UserRole int

const (
    RoleUser UserRole = 1 << iota //bitwise op, iota for auto declaration of const
    RoleAdmin // this would be 10
    RoleShipper // 100
    RoleMod // 1000
)

func (role UserRole) String() string {
    switch role {
    case RoleAdmin:
        return "admin"
    case RoleShipper:
        return "shipper"
    case RoleMod:
        return "mod"
    default:
        return "user"
    }
}

func (role *UserRole) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
			return errors.New(fmt.Sprintf("Failed to unmarshal JSONB value: %v", value))
	}

	var r UserRole
	roleValue := string(bytes)

	if roleValue == "user" {
			r = RoleUser
	} else if roleValue == "admin" {
			r = RoleAdmin
	}

	*role = r
	return nil
}

func (role *UserRole) Value() (driver.Value, error) {
	if role == nil {
			return nil, nil
	}
	return role.String(), nil
}

func (role *UserRole) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", role.String())), nil
}

type User struct {
	common.SQLModel `json:",inline"`
	Email    string   `json:"email" gorm:"column:email;"`
	Password string   `json:"-" gorm:"column:password;"`
	Salt     string   `json:"-" gorm:"column:salt;"`
	LastName string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string  `json:"first_name" gorm:"column:first_name;"`
	Phone    string   `json:"phone" gorm:"column:phone;"`
	Role     UserRole `json:"role" gorm:"column:role;"`
	Status 	 int 			`json:"status" gorm:"column:status;"`
}


//Implementing Requester in const.go
func (u *User) GetUserId() int {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role.String()
}

func (User) TableName() string {
	return "users"
}

type UserCreate struct {
	common.SQLModel `json:",inline"`
	Email    string   `json:"email" gorm:"column:email;"`
	Password string   `json:"password" gorm:"column:password;"`
	LastName string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string  `json:"first_name" gorm:"column:first_name;"`
	Role     string   `json:"-" gorm:"column:role;"`
	Salt     string   `json:"-" gorm:"column:salt;"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

type UserLogin struct {
	Email    string `json:"email" form:"email" gorm:"column:email;"`
	Password string `json:"password" form:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

var (
	ErrEmailOrPasswordInvalid = common.NewCustomError(
			errors.New("email or password invalid"),
			"email or password invalid",
			"ErrUsernameOrPasswordInvalid",
	)

	ErrEmailExisted = common.NewCustomError(
			errors.New("email has already existed"),
			"email has already existed",
			"ErrEmailExisted",
	)
)

