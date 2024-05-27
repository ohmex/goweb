package models

type Domain struct {
	Base
	Name  string  `json:"name" gorm:"type:text"`
	Users []*User `json:"users,omitempty" gorm:"many2many:domain_user;"`
}
