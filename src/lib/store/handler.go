package store

import (
	"errors"
	"log"
	"os"
	"seraphim-cli/lib/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// func createDbTable() {

// }

type StoredConnection struct {
	UpdatedAt int   // Set to current unix seconds on updating or if it is zero on creating
	Updated   int64 `gorm:"autoUpdateTime:nano"` // Use unix nano seconds as updating time
	CreatedAt int64 `gorm:"autoCreateTime"`
	ID        uint  `gorm:"primaryKey"`
	Host      string
	User      string
	SshKey    string
	DefaultDb string
	Provider  string
}

func createStoredConnectionIfNotExists() (*gorm.DB, error) {
	if doesDBExist, storePath := doesStoreDBFileExist(); doesDBExist {
		db, err := gorm.Open(sqlite.Open(storePath), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		db.AutoMigrate(&StoredConnection{})
		return db, nil
	} else {
		log.Fatal("Store database does not exist")
		return nil, errors.New("store db does not exist")
	}

}

// create function to insert one

func doesStoreDBFileExist() (bool, string) {
	storePath, err := config.GetSeraphStore()
	if err == nil {
		if len(storePath) > 0 {
			if _, err := os.Stat(storePath); os.IsNotExist(err) {
				_, err := os.Create(storePath)
				if err != nil {
					log.Fatal(err)
					return false, ""
				}
			}
			return true, storePath
		}
		return false, ""
	} else {
		return false, ""
	}
}

func GetStoredDbConnections() ([]StoredConnection, error) {
	db, err := createStoredConnectionIfNotExists()
	if err != nil {
		log.Fatal(err)
	}
	var storedConnections []StoredConnection
	result := db.Find(&storedConnections)
	if result.RowsAffected > 0 {
		return storedConnections, nil
	} else {
		if err := result.Error; err != nil {
			return nil, err
		} else {
			return nil, errors.New("no results found")
		}
	}

}
