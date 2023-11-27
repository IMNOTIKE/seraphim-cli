package store

import (
	"fmt"
	"seraphim-cli/lib/config"
)

func DoesStoreDBFileExist() {
	storePath, err := config.GetSeraphStore()
	if err == nil {
		fmt.Println(storePath)
	} else {
		fmt.Println(err)
	}
}
