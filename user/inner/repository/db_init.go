package repository

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	dbName := viper.GetString("mysql.dbName")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	charset := viper.GetString("mysql.charset")

	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", dbName, "?charset=" + charset + "&parseTime=true"}, "")
	err := Database(dsn)
	if err != nil {
		panic(err)
	}
}

// 配置mysql
func Database(dsn string) error {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	} //开发模式下就多输出日志
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,  //禁用datatime的精度
		DontSupportRenameIndex:    true,  //不支持重命名索引
		DontSupportRenameColumn:   true,  //不支持重命名列
		SkipInitializeWithVersion: false, //关闭根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("连接db服务失败")
	}
	sqlDB.SetMaxIdleConns(20)                  // 设置空闲连接池
	sqlDB.SetMaxOpenConns(100)                 //设置最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * 30) //最大连接时间
	DB = db
	migration()

	return err
}
