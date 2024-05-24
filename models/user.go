package models

type User struct {
	Base
	Email    string    `json:"email" gorm:"type:varchar(200);"`
	Name     string    `json:"name" gorm:"type:varchar(200);"`
	Password string    `json:"-" gorm:"type:varchar(200);"`
	Tenants  []*Tenant `json:"tenants,omitempty" gorm:"many2many:tenant_user;"`
	Posts    []*Post   `json:"posts,omitempty"`
}
