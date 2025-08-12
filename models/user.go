package models

type Role struct {
	Base
	Name     string `json:"name" gorm:"type:varchar(64);"`
	DomainID uint64
}

type User struct {
	Base
	Email      string    `json:"email" gorm:"type:varchar(200);uniqueIndex"`
	Name       string    `json:"name" gorm:"type:varchar(200);"`
	Password   string    `json:"-" gorm:"type:varchar(200);"`
	Avatar     string    `json:"avatar" gorm:"type:varchar(500);"`
	Provider   string    `json:"provider" gorm:"type:varchar(50);default:'local'"` // local, google, github, etc.
	ProviderID string    `json:"provider_id" gorm:"type:varchar(200);"`            // ID from the OAuth provider
	IsVerified bool      `json:"is_verified" gorm:"default:false"`                 // Email verification status
	Domains    []*Domain `json:"domains,omitempty" gorm:"many2many:domain_users;"`
	Posts      []*Post   `json:"posts,omitempty"`
}

type DomainUser struct {
	DomainID uint64 `gorm:"primaryKey"`
	UserID   uint64 `gorm:"primaryKey"`
	Active   bool   `gorm:"default:false"`
}
