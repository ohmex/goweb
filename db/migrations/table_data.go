package migrations

import (
	"goweb/models"

	"github.com/casbin/casbin/v2"
	ga "github.com/casbin/gorm-adapter/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TableData struct{}

func (TableData) Id() string {
	return "InsertData"
}

func (TableData) Up(db *gorm.DB) {
	adaptor, _ := ga.NewAdapterByDBUseTableName(db, "", "casbin")
	casbin, _ := casbin.NewEnforcer("casbin/model.conf", adaptor)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	system := models.Domain{Name: "System"}
	db.Create(&system)

	role := models.Role{Name: "Admin", DomainID: system.ID}
	db.Create(&role)
	role = models.Role{Name: "Manager", DomainID: system.ID}
	db.Create(&role)
	role = models.Role{Name: "Operator", DomainID: system.ID}
	db.Create(&role)

	reliance := models.Domain{Name: "Reliance"}
	db.Create(&reliance)
	relianceUUID := reliance.UUID.String() //Get UUID port creating table
	casbin.AddPolicy("Admin", relianceUUID, "User", "List")
	casbin.AddPolicy("Admin", relianceUUID, "User", "Read")
	casbin.AddPolicy("Admin", relianceUUID, "User", "Create")
	casbin.AddPolicy("Admin", relianceUUID, "User", "Update")
	casbin.AddPolicy("Admin", relianceUUID, "User", "Delete")
	casbin.AddPolicy("Admin", relianceUUID, "Role", "List")
	casbin.AddPolicy("Admin", relianceUUID, "Role", "Read")
	casbin.AddPolicy("Admin", relianceUUID, "Role", "Create")
	casbin.AddPolicy("Admin", relianceUUID, "Role", "Update")
	casbin.AddPolicy("Admin", relianceUUID, "Role", "Delete")
	casbin.AddPolicy("Admin", relianceUUID, "Post", "List")
	casbin.AddPolicy("Admin", relianceUUID, "Post", "Read")
	casbin.AddPolicy("Admin", relianceUUID, "Post", "Create")
	casbin.AddPolicy("Admin", relianceUUID, "Post", "Update")
	casbin.AddPolicy("Admin", relianceUUID, "Post", "Delete")
	casbin.AddPolicy("Manager", relianceUUID, "User", "List")
	casbin.AddPolicy("Manager", relianceUUID, "User", "Read")
	casbin.AddPolicy("Manager", relianceUUID, "User", "Update")
	casbin.AddPolicy("Manager", relianceUUID, "Role", "List")
	casbin.AddPolicy("Manager", relianceUUID, "Role", "Read")
	casbin.AddPolicy("Operator", relianceUUID, "User", "Read")

	dmart := models.Domain{Name: "DMart"}
	db.Create(&dmart)
	dmartUUID := dmart.UUID.String() //Get UUID port creating table
	casbin.AddPolicy("Admin", dmartUUID, "User", "List")
	casbin.AddPolicy("Admin", dmartUUID, "User", "Read")
	casbin.AddPolicy("Admin", dmartUUID, "User", "Create")
	casbin.AddPolicy("Admin", dmartUUID, "User", "Update")
	casbin.AddPolicy("Admin", dmartUUID, "User", "Delete")
	casbin.AddPolicy("Admin", dmartUUID, "Role", "List")
	casbin.AddPolicy("Admin", dmartUUID, "Role", "Read")
	casbin.AddPolicy("Admin", dmartUUID, "Role", "Create")
	casbin.AddPolicy("Admin", dmartUUID, "Role", "Update")
	casbin.AddPolicy("Admin", dmartUUID, "Role", "Delete")
	casbin.AddPolicy("Admin", dmartUUID, "Post", "List")
	casbin.AddPolicy("Admin", dmartUUID, "Post", "Read")
	casbin.AddPolicy("Admin", dmartUUID, "Post", "Create")
	casbin.AddPolicy("Admin", dmartUUID, "Post", "Update")
	casbin.AddPolicy("Admin", dmartUUID, "Post", "Delete")
	casbin.AddPolicy("Manager", dmartUUID, "User", "List")
	casbin.AddPolicy("Manager", dmartUUID, "User", "Read")
	casbin.AddPolicy("Manager", dmartUUID, "User", "Update")
	casbin.AddPolicy("Manager", dmartUUID, "Role", "List")
	casbin.AddPolicy("Manager", dmartUUID, "Role", "Read")
	casbin.AddPolicy("Operator", dmartUUID, "User", "Read")

	// Setting Super User
	user := models.User{Name: "Sachin", Email: "trulysachin@gmail.com", Password: string(hashedPassword), Domains: []*models.Domain{&reliance, &dmart}}
	db.Create(&user)
	db.Save(&models.DomainUser{UserID: user.ID, DomainID: reliance.ID, Active: true})
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Admin", relianceUUID) // Sachin is Admin of Reliance
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Admin", dmartUUID)    // Sachin is Admin of DMart

	// Setting Users of Reliance
	user = models.User{Name: "UserA", Email: "usera@gmail.com", Password: string(hashedPassword), Domains: []*models.Domain{&reliance}}
	db.Create(&user)
	db.Save(&models.DomainUser{UserID: user.ID, DomainID: reliance.ID, Active: true})
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Admin", relianceUUID)

	user = models.User{Name: "UserB", Email: "userb@gmail.com", Password: string(hashedPassword), Domains: []*models.Domain{&reliance}}
	db.Create(&user)
	db.Save(&models.DomainUser{UserID: user.ID, DomainID: reliance.ID, Active: true})
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Manager", relianceUUID)

	user = models.User{Name: "UserC", Email: "userc@gmail.com", Password: string(hashedPassword), Domains: []*models.Domain{&reliance}}
	db.Create(&user)
	db.Save(&models.DomainUser{UserID: user.ID, DomainID: reliance.ID, Active: true})
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Operator", relianceUUID)

	// Setting Users of DMart
	user = models.User{Name: "UserD", Email: "userd@gmail.com", Password: string(hashedPassword), Domains: []*models.Domain{&dmart}}
	db.Create(&user)
	db.Save(&models.DomainUser{UserID: user.ID, DomainID: dmart.ID, Active: true})
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Admin", dmartUUID)

	user = models.User{Name: "UserE", Email: "usere@gmail.com", Password: string(hashedPassword), Domains: []*models.Domain{&dmart}}
	db.Create(&user)
	db.Save(&models.DomainUser{UserID: user.ID, DomainID: dmart.ID, Active: true})
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Manager", dmartUUID)

	user = models.User{Name: "UserF", Email: "userf@gmail.com", Password: string(hashedPassword), Domains: []*models.Domain{&dmart}}
	db.Create(&user)
	db.Save(&models.DomainUser{UserID: user.ID, DomainID: dmart.ID, Active: true})
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Operator", dmartUUID)
}

func (TableData) Down(db *gorm.DB) {

}
