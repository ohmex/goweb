package migrations

import (
	"goweb/models"

	"gorm.io/gorm"
)

type TableData struct{}

func (TableData) Id() string {
	return "InsertData"
}

func (TableData) Up(db *gorm.DB) {

	user := models.User{Name: "Sachin", Email: "trylysachin@gmail.com", Password: "password"}

	db.Create(&user)

}

func (TableData) Down(db *gorm.DB) {

}
