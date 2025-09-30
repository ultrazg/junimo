package backend

import (
	"fmt"
	"log"

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

func (a *App) ReadConfig() Config {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("解析配置文件失败: %v", err)
	}

	return config
}

func (a *App) UpdateConfig(key string, value any) (bool, string) {
	viper.Set(key, value)

	if err := viper.WriteConfig(); err != nil {
		log.Printf("无法写入配置文件: %v", err)
		return false, "无法写入配置文件"
	}

	log.Println("配置文件已更新")
	return true, "配置文件已更新"
}
