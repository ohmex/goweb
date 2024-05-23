package migrations

import (
	"goweb/models"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TableData struct{}

func (TableData) Id() string {
	return "InsertData"
}

func (TableData) Up(db *gorm.DB) {
	adaptor, _ := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin")
	casbin, _ := casbin.NewEnforcer("casbin/model.conf", adaptor)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	reliance := models.Tenant{Name: "Reliance"}
	db.Create(&reliance)
	relianceUUID := reliance.UUID.String() //Get UUID port creating table
	casbin.AddPolicy("Admin", relianceUUID, "User", "Read")
	casbin.AddPolicy("Admin", relianceUUID, "User", "Create")
	casbin.AddPolicy("Admin", relianceUUID, "User", "Update")
	casbin.AddPolicy("Admin", relianceUUID, "User", "Delete")
	casbin.AddPolicy("Manager", relianceUUID, "User", "Read")
	casbin.AddPolicy("Manager", relianceUUID, "User", "Update")
	casbin.AddPolicy("Operator", relianceUUID, "User", "Read")

	dmart := models.Tenant{Name: "DMart"}
	db.Create(&dmart)
	dmartUUID := dmart.UUID.String() //Get UUID port creating table
	casbin.AddPolicy("Admin", dmartUUID, "User", "Read")
	casbin.AddPolicy("Admin", dmartUUID, "User", "Create")
	casbin.AddPolicy("Admin", dmartUUID, "User", "Update")
	casbin.AddPolicy("Admin", dmartUUID, "User", "Delete")
	casbin.AddPolicy("Manager", dmartUUID, "User", "Read")
	casbin.AddPolicy("Manager", dmartUUID, "User", "Update")
	casbin.AddPolicy("Operator", dmartUUID, "User", "Read")

	// Setting Super User
	user := models.User{Name: "Sachin", Email: "trulysachin@gmail.com", Password: string(hashedPassword), Tenants: []*models.Tenant{&reliance, &dmart}}
	db.Create(&user)
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Admin", relianceUUID) // Sachin is Admin of Reliance
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Admin", dmartUUID)    // Sachin is Admin of DMart

	// Setting Users of Reliance
	user = models.User{Name: "UserA", Email: "usera@gmail.com", Password: string(hashedPassword), Tenants: []*models.Tenant{&reliance}}
	db.Create(&user)
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Admin", relianceUUID)

	user = models.User{Name: "UserB", Email: "userb@gmail.com", Password: string(hashedPassword), Tenants: []*models.Tenant{&reliance}}
	db.Create(&user)
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Manager", relianceUUID)

	user = models.User{Name: "UserC", Email: "userc@gmail.com", Password: string(hashedPassword), Tenants: []*models.Tenant{&reliance}}
	db.Create(&user)
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Operator", relianceUUID)

	// Setting Users of DMart
	user = models.User{Name: "UserD", Email: "userd@gmail.com", Password: string(hashedPassword), Tenants: []*models.Tenant{&dmart}}
	db.Create(&user)
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Admin", dmartUUID)

	user = models.User{Name: "UserE", Email: "usere@gmail.com", Password: string(hashedPassword), Tenants: []*models.Tenant{&dmart}}
	db.Create(&user)
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Manager", dmartUUID)

	user = models.User{Name: "UserF", Email: "userf@gmail.com", Password: string(hashedPassword), Tenants: []*models.Tenant{&dmart}}
	db.Create(&user)
	casbin.AddRoleForUserInDomain(user.UUID.String(), "Operator", dmartUUID)
}

func (TableData) Down(db *gorm.DB) {

}
