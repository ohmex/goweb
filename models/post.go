package models

type Post struct {
	BaseResource
	Title   string `json:"title" gorm:"type:text"`
	Content string `json:"content" gorm:"type:text"`
	UserID  int
	User    User `gorm:"foreignkey:UserID"`
}
