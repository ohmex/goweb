package db

import "gorm.io/gorm"

type Model struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement"`
	Migration string `gorm:"not null"`
	Batch     uint64 `gorm:"not null"`
}

func (Model) TableName() string {
	return "migrations"
}

type Migrator interface {
	Id() string
	Up(db *gorm.DB)
	Down(db *gorm.DB)
}

var (
	setups []Migrator
)

func AddMigrators(setup ...Migrator) {
	setups = append(setups, setup...)
}

func Migrate(db *gorm.DB) (err error) {
	var (
		batch       uint64
		installed   []Model
		installMap  = make(map[string]struct{})
		unInstalled []Migrator
	)

	if err = setupMigrationTable(db); err != nil {
		return err
	}

	if err = db.Order("id asc").Find(&installed).Error; err != nil {
		return err
	}
	for _, v := range installed {
		installMap[v.Migration] = struct{}{}
		if v.Batch > batch {
			batch = v.Batch
		}
	}

	for _, v := range setups {
		if _, ok := installMap[v.Id()]; !ok {
			unInstalled = append(unInstalled, v)
		}
	}
	if err = migrate(db, batch+1, unInstalled...); err != nil {
		return err
	}
	return nil
}

func MigrateUp(db *gorm.DB) (err error) {
	var (
		batch      uint64
		installed  []Model
		installMap = make(map[string]struct{})
	)

	if err = setupMigrationTable(db); err != nil {
		return err
	}

	if err = db.Order("id asc").Find(&installed).Error; err != nil {
		return err
	}
	for _, v := range installed {
		installMap[v.Migration] = struct{}{}
		if v.Batch > batch {
			batch = v.Batch
		}
	}
	for _, v := range setups {
		if _, ok := installMap[v.Id()]; !ok {
			return migrate(db, batch+1, v)
		}
	}
	return nil
}

func MigrateDown(db *gorm.DB) (err error) {
	var (
		batch     uint64
		installed []Model
		downIdMap = make(map[string]struct{})
	)

	if err = setupMigrationTable(db); err != nil {
		return err
	}

	if err = db.Order("batch desc, id desc").Find(&installed).Error; err != nil {
		return err
	}
	for _, v := range installed {
		if batch == 0 {
			batch = v.Batch
		}
		if v.Batch == batch {
			downIdMap[v.Migration] = struct{}{}
		}
	}
	for i := len(setups) - 1; i >= 0; i-- {
		if _, ok := downIdMap[setups[i].Id()]; !ok {
			continue
		}
		setups[i].Down(db)
	}
	return db.Where("batch=?", batch).Delete(&Model{}).Error
}

func migrate(db *gorm.DB, batch uint64, migrators ...Migrator) (err error) {
	for _, v := range migrators {
		v.Up(db)
		if err = db.Create(&Model{
			Migration: v.Id(),
			Batch:     batch,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func setupMigrationTable(db *gorm.DB) error {
	model := Model{}
	if !db.Migrator().HasTable(&model) {
		if err := db.Migrator().CreateTable(&model); err != nil {
			return err
		}
	}
	return nil
}
