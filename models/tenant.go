package models

type Tenant struct {
	Base
	Name  string  `json:"name" gorm:"type:text"`
	Users []*User `json:"users" gorm:"many2many:tenant_user;"`
}
