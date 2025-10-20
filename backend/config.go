package backend

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		viper.SetDefault("theme", "light")
		viper.SetDefault("game_path", "")
		viper.SetDefault("backup_path", "")
		viper.SetDefault("disabled_path", "")
		viper.SetDefault("mods", []ModManifestJson{})
		viper.SetDefault("smapi_version", "")
		viper.SetDefault("disabled", false)

		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("读取配置文件失败: %v\n", err)
			if err = viper.SafeWriteConfigAs("config.json"); err != nil {
				fmt.Printf("写入配置文件失败: %v\n", err)
			}
		}
	} else {
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("读取配置文件失败: %v\n", err)
		}
	}
}

func (a *App) ReadConfig(key string) any {
	return viper.Get(key)
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
