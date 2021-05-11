package config

import (
	"errors"
	"fmt"
	"github.com/guzhongzhi/gmicro/logger"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

const (
	EnvDev  = "dev"
	EnvQA   = "qa"
	EnvProd = "prod"
)

func LoadConfigByFiles(path, env string, bootstrap interface{}, logger logger.SuperLogger) (error) {

	err := readConfigFiles(env, path, bootstrap, logger)
	return err
}

func generateCfgKeys(t reflect.Type, parentPath string) map[string]reflect.Type {
	names := make(map[string]reflect.Type)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return names
	}

	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		fieldName := t.Field(i).Name
		fullFieldName := parentPath
		if fullFieldName != "" {
			fullFieldName += "/" + fieldName
		} else {
			fullFieldName = fieldName
		}

		st := t.Field(i).Type
		if st.Kind() == reflect.Ptr {
			st = st.Elem()
		}
		if st.Kind() == reflect.Struct {
			subNames := generateCfgKeys(st, fullFieldName)
			for key, t := range subNames {
				names[key] = t
			}
		} else {
			names[fullFieldName] = st
		}
	}
	return names
}

func readConfigFiles(env, dir string, out interface{}, logger logger.SuperLogger) error {
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

	t := reflect.TypeOf(out)
	keys := generateCfgKeys(t, "")

	for key, t := range keys {
		envName := strings.Replace(strings.ToUpper(key), "/", "_", -1)
		logger.Debug(envName, " ", t.Kind(), " ", t.Name())
		if os.Getenv(envName) == "" {
			continue
		}
		switch viper.Get(key).(type) {
		case []string:
			viper.Set(key, []string{os.Getenv(envName)})
		case string:
			viper.Set(key, os.Getenv(envName))
		default:
			panic(fmt.Errorf("unsupport ENV data type '%v' for the field of '%s'", t.Kind(), envName))
		}
	}
	return viper.Unmarshal(out)
}
