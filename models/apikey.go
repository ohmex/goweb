package models

import "time"

type APIKey struct {
	Base
	UserID       uint64
	ApiKey       string    `json:"api_key" gorm:"type:varchar(255);"`
	HashedSecret string    `json:"-" gorm:"type:varchar(255);"`
	Active       bool      `gorm:"default:false"`
	UsedAt       time.Time `json:"used_at" gorm:"default:null"`
}

// "User: ":"usera@gmail.com","Key: ":"dc765138784476a8da83f8de90195b9e","Secret: ":"90aaa42ea05f9a9199d1333bd263825dd53d0ab5724a4f8f804c05b5997d0a4b"
