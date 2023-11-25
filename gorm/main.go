package main

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Device struct {
	Uuid      uuid.UUID `gorm:"type:uuid;primaryKey;not null;default:uuid_generate_v4()"`
	Name      string
	Arch      string
	Os        string
	Vulkan    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func main() {
	dsn := "host=xxx user=xxx password=xxx dbname=backend sslmode=disable TimeZone=Europe/Prague"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err2 := db.AutoMigrate(&Device{})
	if err2 != nil {
		return
	}

	// Create
	db.Create(&Device{Name: "PC", Os: "linux"})

	// Read
	var device Device
	//db.First(&device, 1)                // find device with integer primary key
	db.First(&device, "name = ?", "PC") // find device with code D42

	// Update - update device's price to 200
	//db.Model(&device).Update("Price", 200)
	// Update - update multiple fields
	//db.Model(&device).Updates(Device{Name: "notebook", Os: "windows"}) // non-zero fields
	//db.Model(&device).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete device
	//db.Delete(&device)

	// return all columns
	var devices []Device
	db.Find(&devices)
	fmt.Printf("%v", devices)
}
