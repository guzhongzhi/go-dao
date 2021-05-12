package config

import (
	"errors"
	"fmt"
	"github.com/guzhongzhi/gmicro/logger"
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
	viper.AutomaticEnv()
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

	allSettings := viper.AllSettings()
	for key, _ := range keys {
		viperKey := strings.Replace(strings.ToLower(key), "/", ".", -1)
		if _, ok := allSettings[viperKey]; !ok {
			viper.Set(viperKey, viper.Get(viperKey))
		}
	}
	return viper.Unmarshal(out)
}
