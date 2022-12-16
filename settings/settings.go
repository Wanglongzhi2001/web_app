package settings

import (
	"flag"
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func Init() (err error) {
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")

	var filename string
	flag.StringVar(&filename, "filename", "config.yaml", "配置文件")
	flag.Parse()
	viper.SetConfigFile(filename)
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return
}
