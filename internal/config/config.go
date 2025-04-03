package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Providers map[string]ProviderConfig `mapstructure:"providers"`
	Deployers map[string]DeployerConfig `mapstructure:"deployers"`
}

type ProviderConfig struct {
	Type     string            `mapstructure:"type"`
	Settings map[string]string `mapstructure:"settings"`
}

type DeployerConfig struct {
	Type     string            `mapstructure:"type"`
	Settings map[string]string `mapstructure:"settings"`
}

var globalConfig *Config

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 创建配置目录
	configDir := filepath.Join(os.Getenv("HOME"), ".fuckssl")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}
	viper.AddConfigPath(configDir)

	// 设置默认值
	viper.SetDefault("providers", make(map[string]ProviderConfig))
	viper.SetDefault("deployers", make(map[string]DeployerConfig))

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，创建默认配置
			configPath := filepath.Join(configDir, "config.yaml")
			if err := viper.WriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to create default config: %v", err)
			}
			fmt.Printf("Created default config file at %s\n", configPath)
		} else {
			return fmt.Errorf("failed to read config: %v", err)
		}
	}

	globalConfig = &Config{}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return nil
}

func GetConfig() *Config {
	return globalConfig
}

func SaveConfig() error {
	viper.Set("providers", globalConfig.Providers)
	viper.Set("deployers", globalConfig.Deployers)
	return viper.WriteConfig()
}
