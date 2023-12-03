package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

func Init(configFileName string) (err error) {
	viper.SetConfigFile(configFileName)
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config read:%s\n", err))
		return err
	}
	err = viper.Unmarshal(Conf)
	if err != nil {
		panic(fmt.Errorf("Fatal error config unmarshal:%s\n", err))
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		err = viper.Unmarshal(Conf)
		if err != nil {
			panic(fmt.Errorf("Fatal error config unmarshal:%s\n", err))
		}
	})
	return nil
}

type AppConfig struct {
	Name            string `mapstructure:"name"`
	Mode            string `mapstructure:"mode"`
	StartTime       string `mapstructure:"start_time"`
	Version         string `mapstructure:"version"`
	MachineId       int64  `mapstructure:"machine_id"`
	Port            int    `mapstructure:"port"`
	*LogConfig      `mapstructure:"log"`
	*MySQLConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*RabbitMQConfig `mapstructure:"rabbitmq"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxSize    int    `mapstructure:"max_size"`
}

type MySQLConfig struct {
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DBName      string `mapstructure:"dbname"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
}

type RedisConfig struct {
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	DB       int    `mapstructure:"db"`
	Port     int    `mapstructure:"port"`
	PoolSize int    `mapstructure:"pool_size"`
}

type RabbitMQConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}
