package models

type User struct {
	Base
	Email    string    `json:"email" gorm:"type:varchar(200);"`
	Name     string    `json:"name" gorm:"type:varchar(200);"`
	Password string    `json:"-" gorm:"type:varchar(200);"`
	Domains  []*Domain `json:"domains,omitempty" gorm:"many2many:domain_user;"`
	Posts    []*Post   `json:"posts,omitempty"`
}
