package models

type Post struct {
	BaseResource
	Title   string `json:"title" gorm:"type:text"`
	Content string `json:"content" gorm:"type:text"`
	UserID  uint
	User    User `gorm:"foreignkey:UserID"`
}
