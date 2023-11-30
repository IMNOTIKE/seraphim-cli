package config

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type BrandingConfig struct {
	Name string `mapstructure:"name"`
}

type StoredConnection struct {
	Host           string `mapstructure:"host"`
	User           string `mapstructure:"user"`
	Port           int    `mapstructure:"port"`
	Password       string `mapstructure:"pwd"`
	Provider       string `mapstructure:"provider"`
	DefaltDatabase string `mapstructure:"default_database"`
}

type SeraphimConfig struct {
	Version           string                        `mapstructure:"version"`
	BrandingConfig    BrandingConfig                `mapstructure:"branding"`
	StoredConnections []map[string]StoredConnection `mapstructure:"stored_connections"`
}

func RemoveStoredConnection(key string, index int, isInArray bool) tea.Msg {
	if isInArray {
		delete(viper.Get((key + "[" + string(index) + "]")).(map[string]interface{}), "key")
	}
	delete(viper.Get(key).(map[string]interface{}), "key")
	return tea.ExitAltScreen
}
