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

// decoders that used while Unmarshal viper settings to config object
var decoders []viper.DecoderConfigOption

//register decoder
func RegisterDecoder(option viper.DecoderConfigOption) {
	decoders = append(decoders, option)
}

func LoadConfigFiles(path, env string, bootstrap interface{}, logger logger.SuperLogger, envPrefix string) (error) {
	err := readConfigFiles(env, path, bootstrap, logger, envPrefix)
	return err
}

func GenerateCfgKeys(t reflect.Type, parentPath string) map[string]reflect.Type {
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
			subNames := GenerateCfgKeys(st, fullFieldName)
			for key, t := range subNames {
				names[key] = t
			}
		} else {
			names[fullFieldName] = st
		}
	}
	return names
}

func readConfigFiles(env, dir string, out interface{}, logger logger.SuperLogger, envPrefix string) error {
	if env != EnvDev && env != EnvQA && env != EnvProd {
		return errors.New("invalid env param")
	}
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetConfigFile(dir + "/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("load config by viper failed, err=%w", err)
	}

	envFile := fmt.Sprintf("%s/config.%s.yaml", dir, env)
	f, err := os.Open(envFile)
	if err != nil {
		logger.Warnf("no env config file: '%s'", envFile)
	} else {
		defer f.Close()
		viper.MergeConfig(f)
	}

	allSettings := viper.AllSettings()

	t := reflect.TypeOf(out)
	keys := GenerateCfgKeys(t, "")
	for key, _ := range keys {
		viperKey := strings.Replace(strings.ToLower(key), "/", ".", -1)
		if _, ok := allSettings[viperKey]; !ok {
			viper.Set(viperKey, viper.Get(viperKey))
		}
	}
	return viper.Unmarshal(out, decoders...)
}
