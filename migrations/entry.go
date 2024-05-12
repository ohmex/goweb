package main

import (
	"goweb/migrations/list"

	gmEngine "github.com/ShkrutDenis/go-migrations/engine"
	gmStore "github.com/ShkrutDenis/go-migrations/engine/store"
)

func main() {
	e := gmEngine.NewEngine()
	e.Run(getMigrationsList())
	e.GetConnector().Close()
}

func getMigrationsList() []gmStore.Migratable {
	return []gmStore.Migratable{
		&list.CreateUserTable{},
		&list.CreatePostTable{},
		&list.UpdateUserTable{},
	}
}
