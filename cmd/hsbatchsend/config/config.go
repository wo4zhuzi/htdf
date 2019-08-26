package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	globalConfig   *viper.Viper
	onceInitConfig sync.Once
)

func initConfig() {
	globalConfig = NewConfig()
}

func NewConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")

	return v
}

func GetInstanceConfig() *viper.Viper {
	onceInitConfig.Do(func() {
		initConfig()
	})
	return globalConfig
}

func Init(strFileName string, strPath string) error {
	v := GetInstanceConfig()
	v.SetConfigName(strFileName)
	v.AddConfigPath(strPath)

	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("read config file error|err=%s|fileName=%s|strPath=%s\n", err, strFileName, strPath)
		return err
	}

	return nil
}

func relativePath(basedir string, path *string) {
	p := *path
	if p != "" && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetInt(s string, defaultValue int) int {
	v := GetInstanceConfig().GetString(s)

	if v != "" {
		if r, err := strconv.Atoi(v); err == nil {
			return r
		}

	}
	return defaultValue

}

func GetInt64(s string, defaultValue int64) int64 {
	v := GetInstanceConfig().GetString(s)

	if v != "" {
		if r, err := strconv.ParseInt(v, 10, 64); err == nil {
			return r
		}

	}
	return defaultValue

}

func GetString(s, defaultValue string) string {
	v := GetInstanceConfig().GetString(s)

	if v == "" {
		return defaultValue

	}
	return v

}
