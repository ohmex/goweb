package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uint      `gorm:"primarykey"`
	UUID      uuid.UUID `gorm:"type:char(36);index;unique;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.UUID = uuid.New()
	return nil
}

// All resource other than user should use this.
type BaseResource struct {
	Base
	Tenant uuid.UUID `gorm:"type:char(36);"`
}
