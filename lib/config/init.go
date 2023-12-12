package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/coreybutler/go-fsutil"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	pathInputTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
)

type InitResult struct {
	Err error
	Msg string
}

func InitConfig() InitResult {

	file := viper.ConfigFileUsed()
	//path := strings.Join(strings.Split(file, "/")[], "/")

	if val, err := os.Stat(file); err == nil {
		fmt.Println(val)
		return createDefaultConfigFile(file)
	} else if os.(err) {
		return InitResult{
			Err: err,
			Msg: pathInputTitleStyle.Render("Default configuration file already exist, you can reset it to default values by running\n'seraphim config reset'"),
		}
	}

}

func createDefaultConfigFile(file string) InitResult {

	content, err := yaml.Marshal(SeraphimConfig{})
	if err != nil {
		return InitResult{
			Err: err,
			Msg: "",
		}
	}

	writeError := fsutil.WriteTextFile(file, string(content))
	if writeError != nil {
		return InitResult{
			Err: err,
			Msg: "",
		}
	} else {
		return InitResult{
			Err: nil,
			Msg: pathInputTitleStyle.Render(fmt.Sprintf("Successfully created default configuration file at: %s", file)),
		}
	}

}
