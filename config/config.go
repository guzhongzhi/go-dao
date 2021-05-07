package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	EnvDev  = "dev"
	EnvQA   = "qa"
	EnvProd = "prod"
)

func LoadConfigByFiles(path, env string, bootstrap interface{}) (error) {
	configDir := filepath.Join(filepath.Dir(filepath.Dir(os.Args[0])), "configs")

	err := readConfigFiles(env, configDir, bootstrap)
	return err
}

func readConfigFiles(env, dir string, out interface{}) error {
	if env != EnvDev && env != EnvQA && env != EnvProd {
		return errors.New("invalid env param")
	}

	viper.SetConfigFile(dir + "/config.prod.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("load config by viper failed, err=%w", err)
	}

	if env != EnvProd {
		envViper := viper.New()
		envViper.AddConfigPath(dir)
		envViper.SetConfigName("config." + env)
		if err := envViper.ReadInConfig(); err != nil {
			return fmt.Errorf("load config by viper failed, err=%w", err)
		}
		viper.MergeConfigMap(envViper.AllSettings())
	}

	for _, key := range viper.AllKeys() {
		envName := strings.Replace(strings.ToUpper(key), ".", "_", -1)
		if os.Getenv(envName) == "" {
			continue
		}
		switch viper.Get(key).(type) {
		case []string:
			viper.Set(key, []string{os.Getenv(envName)})
		case string:
			viper.Set(key, os.Getenv(envName))
		default:
			fmt.Println("unsupport env value data type")
		}
	}
	return viper.Unmarshal(out)
}
