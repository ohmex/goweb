package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uint64         `json:"-" gorm:"primarykey"`
	UUID      uuid.UUID      `json:"uuid" gorm:"type:char(36);index;unique;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.UUID = uuid.New()
	return nil
}

// All resource other than user should use this.
type BaseResource struct {
	Base
	Domain uuid.UUID `json:"domain" gorm:"type:char(36);"`
}
