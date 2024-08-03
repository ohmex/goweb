package models

type Role struct {
	Base
	Name     string `json:"name" gorm:"type:varchar(64);"`
	DomainID uint64
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
	DomainID uint64 `gorm:"primaryKey"`
	UserID   uint64 `gorm:"primaryKey"`
	Active   bool   `gorm:"default:false"`
}
