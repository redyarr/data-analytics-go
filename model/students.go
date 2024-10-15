package model

import "gorm.io/gorm"

type Student struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Age    int
	Grade  float32
	Gender string
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Student{})

}
