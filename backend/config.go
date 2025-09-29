package backend

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.SetDefault("theme", "light")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		if err = viper.SafeWriteConfigAs("config.json"); err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
		}
	}
}
