package config

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/coreybutler/go-fsutil"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	pathInputTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
)

type ConfigOperationResult struct {
	Err error
	Msg string
}

type BrandingConfig struct {
	Name string `mapstructure:"name"`
}

type StoredConnection struct {
	Host            string `mapstructure:"host"`
	User            string `mapstructure:"user"`
	Port            int    `mapstructure:"port"`
	Password        string `mapstructure:"password"`
	Provider        string `mapstructure:"provider"`
	DefaultDatabase string `mapstructure:"default_database"`
}

type SeraphimConfig struct {
	Version            string                        `mapstructure:"version"`
	Branding           BrandingConfig                `mapstructure:"branding"`
	Stored_Connections []map[string]StoredConnection `mapstructure:"stored_connections"`
	Default_dump_path  string                        `mapstructure:"default_dump_path"`
}

func AddConnection(withConf bool, conf SeraphimConfig, newConn StoredConnection, tag string) ConfigOperationResult {

	file := viper.ConfigFileUsed()
	var config SeraphimConfig
	if withConf {
		config = conf
	} else {
		viper.Unmarshal(&config)
	}
	formattedTag := strings.Replace(strings.Trim(tag, " "), " ", "_", -1)
	newConnMapped := make(map[string]StoredConnection)
	newConnMapped[formattedTag] = newConn
	updatedConnections := append(config.Stored_Connections, newConnMapped)

	config.Stored_Connections = updatedConnections

	content, err := yaml.Marshal(config)
	if err != nil {
		return ConfigOperationResult{
			Err: err,
			Msg: "",
		}
	}

	writeError := fsutil.WriteTextFile(file, string(content))
	if writeError != nil {
		return ConfigOperationResult{
			Err: err,
			Msg: "",
		}
	} else {
		return ConfigOperationResult{
			Err: nil,
			Msg: "Successfully Added connection",
		}
	}
}

func EditConnection(conf SeraphimConfig, oldConn StoredConnection, newConn StoredConnection, oldTag string, newTag string) ConfigOperationResult {

	file := viper.ConfigFileUsed()
	if oldTag == newTag {
		// Only edit it
		if oldConn == newConn {
			return ConfigOperationResult{
				Err: nil,
				Msg: "Nothing to edit",
			}
		}
		anyMatch := false
		for i, v := range conf.Stored_Connections {
			for tag, conn := range v {
				if tag == oldTag {
					conn = newConn
					conf.Stored_Connections[i][tag] = conn
					anyMatch = true
					break
				}
			}
		}
		if anyMatch {
			content, err := yaml.Marshal(conf)
			if err != nil {
				return ConfigOperationResult{
					Err: err,
					Msg: "",
				}
			}

			writeError := fsutil.WriteTextFile(file, string(content))
			if writeError != nil {
				return ConfigOperationResult{
					Err: err,
					Msg: "",
				}
			} else {
				return ConfigOperationResult{
					Err: nil,
					Msg: "Successfully edited selected connection",
				}
			}
		}

	} else {
		// Remove old one and insert new one
		for i, v := range conf.Stored_Connections {
			for tag := range v {
				if tag == oldTag {
					conf.Stored_Connections = append(conf.Stored_Connections[:i], conf.Stored_Connections[i+1:]...)
				}
			}
		}
		if addResult := AddConnection(true, conf, newConn, newTag); addResult.Err != nil {
			return addResult
		}
	}

	return ConfigOperationResult{
		Err: nil,
		Msg: "Successfully edited selected connection",
	}
}

func InitConfig() ConfigOperationResult {

	file := viper.ConfigFileUsed()

	if _, err := os.Stat(file); err != nil {
		return createDefaultConfigFile()
	} else {
		return ConfigOperationResult{
			Err: err,
			Msg: pathInputTitleStyle.Render("Default configuration file already exist, you can reset it to default values by running\n'seraphim config reset'"),
		}
	}

}

func createDefaultConfigFile() ConfigOperationResult {

	homePath, pathErr := os.UserHomeDir()
	if pathErr != nil {
		return ConfigOperationResult{
			Err: pathErr,
			Msg: "",
		}
	}

	content, err := yaml.Marshal(SeraphimConfig{})
	if err != nil {
		return ConfigOperationResult{
			Err: err,
			Msg: "",
		}
	}

	path := fmt.Sprintf("%s/%s/%s", homePath, ".config", "seraphim")
	mkDirerr := os.MkdirAll(path, os.ModePerm)
	if mkDirerr != nil {
		return ConfigOperationResult{
			Err: mkDirerr,
			Msg: "",
		}
	}

	filePath := fmt.Sprintf("%s/%s", path, "seraphim.yaml")
	writeError := fsutil.WriteTextFile(filePath, string(content))
	if writeError != nil {
		return ConfigOperationResult{
			Err: err,
			Msg: "",
		}
	} else {
		return ConfigOperationResult{
			Err: nil,
			Msg: pathInputTitleStyle.Render(fmt.Sprintf("Successfully created default configuration file at: %s", filePath)),
		}
	}

}

func ResetConfig() ConfigOperationResult {

	configFilePath := viper.ConfigFileUsed()

	content, err := yaml.Marshal(SeraphimConfig{})
	if err != nil {
		return ConfigOperationResult{
			Err: err,
			Msg: "",
		}
	}

	writeError := fsutil.WriteTextFile(configFilePath, string(content))
	if writeError != nil {
		return ConfigOperationResult{
			Err: err,
			Msg: "",
		}
	} else {
		return ConfigOperationResult{
			Err: nil,
			Msg: pathInputTitleStyle.Render("Successfully resetted configuration file to default values"),
		}
	}

}

func SaveConfig(conf SeraphimConfig) ConfigOperationResult {
	
	file := viper.ConfigFileUsed()

	content, err := yaml.Marshal(conf)
	if err != nil {
		return ConfigOperationResult{
			Err: err,
			Msg: "",
		}
	}

	writeError := fsutil.WriteTextFile(file, string(content))
	if writeError != nil {
		return ConfigOperationResult{
			Err: err,
			Msg: "",
		}
	} else {
		return ConfigOperationResult{
			Err: nil,
			Msg: "Successfully Saved connection",
		}
	}

}

func ClearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear") // Linux or MacOS
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") // Windows
	default:
		fmt.Println("Unsupported operating system")
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}