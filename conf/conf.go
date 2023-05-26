package conf

import (
	"os"

	"github.com/spf13/viper"
)

func InitConfig() {
	workDir, _ := os.Getwd() //获取工作目录
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/ToDoList_grpc/conf/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}
