package settings

import (
	"fmt"

	dotconfig "github.com/DeanPDX/dotconfig"
)

type SettingsT struct {
	Port string `env:"PORT"`
}

var Settings SettingsT

func init() {
	var err error
	Settings, err = dotconfig.FromFileName[SettingsT](".env")
	if err != nil {
		fmt.Printf("Error: %v.", err)
	}
}
