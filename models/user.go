package models

type Role struct {
	Base
	Name     string `json:"name" gorm:"type:varchar(64);"`
	DomainID int
}

type User struct {
	Base
	Email    string    `json:"email" gorm:"type:varchar(200);"`
	Name     string    `json:"name" gorm:"type:varchar(200);"`
	Password string    `json:"-" gorm:"type:varchar(200);"`
	Domains  []*Domain `json:"domains,omitempty" gorm:"many2many:domain_users;"`
	Posts    []*Post   `json:"posts,omitempty"`
}

type DomainUser struct {
	DomainID int  `gorm:"primaryKey"`
	UserID   int  `gorm:"primaryKey"`
	Active   bool `gorm:"default:false"`
}
