package config

import (
	"github.com/spf13/viper"
	"os"
)

var Conf *Config

type Config struct {
	System *System `yaml:"system"`
	Mysql  *Mysql  `yaml:"mysql"`
	Redis  *Redis  `yaml:"redis"`
	Paths  *Paths  `yaml:"paths"`
	Aliyun *Aliyun `yaml:"aliyun"`
}

type System struct {
	OS           string `yaml:"os"`
	Status       string `yaml:"status"`
	WorkerID     int64  `yaml:"worker_id" mapstructure:"worker_id"`
	DataCenterID int64  `yaml:"data_center_id" mapstructure:"data_center_id"`
}

type Mysql struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type Redis struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Database int    `yaml:"database" mapstructure:"database"`
	Network  string `yaml:"network" mapstructure:"network"`
	Password string `yaml:"password" mapstructure:"password"`
}

type Paths struct {
	DefaultAvatarPath string `yaml:"default_avatar_path" mapstructure:"default_avatar_path"`
}

type Aliyun struct {
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret" mapstructure:"access_key_secret"`
	Endpoint        string `yaml:"endpoint" mapstructure:"endpoint"`
}

// InitConfig initializes the configuration for the project
// and unmarshall the configuration into the global variable "Conf"
func InitConfig() {
	wd, _ := os.Getwd()

	configDIR := os.Getenv("CONFIG_DIR")
	viper.AddConfigPath(configDIR) // auto

	viper.AddConfigPath(wd + "/config/")   // linux
	viper.AddConfigPath(wd + "\\config\\") // windows
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
}
