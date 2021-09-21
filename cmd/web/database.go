package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDb(filename string) (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, err
	}

	if true {
		fmt.Println("Start migration")
		err = db.AutoMigrate(
			&Todo{},
			&User{},
		)
		if err != nil {
			return nil, err
		}
		fmt.Println("End migration")
	} else {
		fmt.Println("No migration")
	}

	return db, err
}
