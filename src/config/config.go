package config

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coreybutler/go-fsutil"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type BrandingConfig struct {
	Name string `mapstructure:"name"`
}

type StoredConnection struct {
	Host           string `mapstructure:"host"`
	User           string `mapstructure:"user"`
	Port           int    `mapstructure:"port"`
	Password       string `mapstructure:"password"`
	Provider       string `mapstructure:"provider"`
	DefaltDatabase string `mapstructure:"defaltdatabase"`
}

type SeraphimConfig struct {
	Version            string                        `mapstructure:"version"`
	Branding           BrandingConfig                `mapstructure:"branding"`
	Stored_Connections []map[string]StoredConnection `mapstructure:"stored_connections"`
	Default_dump_path  string                        `mapstructure:"default_dump_path"`
}

// SHOULD ASK FOR CONFIRMATION
func RemoveStoredConnection(index int, keystoremove ...string) tea.Msg {

	file := viper.ConfigFileUsed()
	var config SeraphimConfig
	viper.Unmarshal(&config)

	for _, key := range keystoremove {
		delete(config.Stored_Connections[index], key)
	}

	if index != -1 {
		config.Stored_Connections = append(config.Stored_Connections[:index], config.Stored_Connections[index+1:]...)
	}

	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fsutil.WriteTextFile(file, string(content))

	return "Removed config entry"
}

func AddConnection(newConn StoredConnection, tag string) tea.Msg {

	file := viper.ConfigFileUsed()
	var config SeraphimConfig
	viper.Unmarshal(&config)
	formattedTag := strings.Replace(strings.Trim(tag, " "), " ", "_", -1)
	newConnMapped := make(map[string]StoredConnection)
	newConnMapped[formattedTag] = newConn
	updatedConnections := append(config.Stored_Connections, newConnMapped)

	config.Stored_Connections = updatedConnections

	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fsutil.WriteTextFile(file, string(content))

	return "Added new stored connection"
}
