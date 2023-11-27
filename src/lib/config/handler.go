package config

import (
	"errors"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type StoreConfig struct {
	STORE_DIR string
}

type BrandingConfig struct {
	Name string
}

type SeraphConfig struct {
	Version        string
	StoreConfig    StoreConfig    `toml:"store"`
	BrandingConfig BrandingConfig `toml:"branding"`
}

const CONFIG_PATH = "/.config/seraphim/settings.toml"

func GetSerpahConfig() SeraphConfig {
	var homePath string
	if v, err := os.UserHomeDir(); err == nil {
		homePath = v
	} else {
		log.Fatal(err)
	}
	fullPath := homePath + "" + CONFIG_PATH
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		_, err := os.Create(fullPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	var conf SeraphConfig
	confData, readError := os.ReadFile(fullPath)
	if readError != nil {
		log.Fatal(readError)
	}
	if _, err := toml.Decode(string(confData), &conf); err != nil {
		log.Fatal(readError)
	}
	return conf
}

func GetSeraphStore() (string, error) {
	var message string
	conf := GetSerpahConfig()
	if conf == (SeraphConfig{}) {
		message = "Couldn't get configs"
		log.Fatal(message)
		return "", errors.New(message)
	} else {
		if path := conf.StoreConfig.STORE_DIR; len(path) > 0 {
			return path, nil
		} else {
			message = "No path for store was set"
			log.Fatal(message)
		}
		return "", errors.New(message)
	}

}
