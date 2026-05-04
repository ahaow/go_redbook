package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app" yaml:"app"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Redis    RedisConfig    `mapstructure:"redis" yaml:"redis"`
	Jwt      JwtConfig      `mapstructure:"jwt" yaml:"jwt"`
	Upload   UploadConfig   `mapstructure:"upload" yaml:"upload"`
}

type AppConfig struct {
	Name string `mapstructure:"name" yaml:"name"`
	Port string `mapstructure:"port" yaml:"port"`
}

type DatabaseConfig struct {
	Mode         string `mapstructure:"mode" yaml:"mode"`
	Dsn          string `mapstructure:"dsn" yaml:"dsn"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenCons  int    `mapstructure:"max_open_cons" yaml:"max_open_cons"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address" yaml:"address"`
	Password string `mapstructure:"password" yaml:"password"`
	Db       int    `mapstructure:"db" yaml:"db"`
}

type JwtConfig struct {
	Expires int    `mapstructure:"expires" yaml:"expires"`
	Issuer  string `mapstructure:"issuer" yaml:"issuer"`
	Key     string `mapstructure:"key" yaml:"key"`
}

type UploadConfig struct {
	Size int64  `mapstructure:"size" yaml:"size"`
	Dir  string `mapstructure:"dir" yaml:"dir"`
}

func InitConfig() (*Config, error) {
	configName := "config"

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置失败 [%s.yaml]: %w", configName, err)
	}
	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("配置文件解析失败: %w", err)
	}
	fmt.Println("配置加载成功:", viper.ConfigFileUsed())
	return cfg, nil
}
