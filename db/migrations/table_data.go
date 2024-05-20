package migrations

import (
	"goweb/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TableData struct{}

func (TableData) Id() string {
	return "InsertData"
}

func (TableData) Up(db *gorm.DB) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user := models.User{Name: "Sachin", Email: "trulysachin@gmail.com", Password: string(hashedPassword)}
	db.Create(&user)
	user = models.User{Name: "UserA", Email: "usera@gmail.com", Password: string(hashedPassword)}
	db.Create(&user)
	user = models.User{Name: "UserB", Email: "userb@gmail.com", Password: string(hashedPassword)}
	db.Create(&user)
	user = models.User{Name: "UserC", Email: "userc@gmail.com", Password: string(hashedPassword)}
	db.Create(&user)
	user = models.User{Name: "UserD", Email: "userd@gmail.com", Password: string(hashedPassword)}
	db.Create(&user)
	user = models.User{Name: "UserE", Email: "usere@gmail.com", Password: string(hashedPassword)}
	db.Create(&user)
	user = models.User{Name: "UserF", Email: "userf@gmail.com", Password: string(hashedPassword)}
	db.Create(&user)
}

func (TableData) Down(db *gorm.DB) {

}
