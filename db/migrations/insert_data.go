package migrations

import (
	"goweb/models"

	"gorm.io/gorm"
)

type InsertData struct{}

func (InsertData) Id() string {
	return "InsertData"
}

func (InsertData) Up(db *gorm.DB) {

	user := models.User{Name: "Sachin", Email: "trylysachin@gmail.com", Password: "password"}

	db.Create(&user)

}

func (InsertData) Down(db *gorm.DB) {

}
