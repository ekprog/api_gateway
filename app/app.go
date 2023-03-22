package app

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

func InitApp(rootDir ...string) error {

	basePath := "."
	if len(rootDir) != 0 {
		basePath = rootDir[0]
	}

	yamlPath := "config"
	yamlApp := "app.yaml"
	yamlFullPath := path.Join(basePath, yamlPath, yamlApp)
	if _, err := os.Stat(yamlFullPath); errors.Is(err, os.ErrNotExist) {
		return errors.Errorf("create config/app.yaml first in root directory")
	}

	//viper.AddConfigPath(yamlPath)
	viper.SetConfigFile(yamlFullPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "viper cannot read config")
	}

	envPath := path.Join(basePath, ".env")
	if _, err := os.Stat(envPath); err == nil {
		err := godotenv.Load(envPath)
		if err != nil {
			return errors.Wrap(err, "godotenv cannot load config")
		}
	}

	return nil
}
