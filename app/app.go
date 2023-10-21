package app

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

func InitApp(rootDir ...string) (context.Context, context.CancelFunc, error) {

	// Graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	basePath := "."
	if len(rootDir) != 0 {
		basePath = rootDir[0]
	}

	yamlPath := "config"
	yamlApp := "app.yaml"
	yamlFullPath := path.Join(basePath, yamlPath, yamlApp)
	if _, err := os.Stat(yamlFullPath); errors.Is(err, os.ErrNotExist) {
		return nil, nil, errors.Errorf("create config/app.yaml first in root directory")
	}

	//viper.AddConfigPath(yamlPath)
	viper.SetConfigFile(yamlFullPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "viper cannot read config")
	}

	envPath := path.Join(basePath, ".env")
	if _, err := os.Stat(envPath); err == nil {
		err := godotenv.Load(envPath)
		if err != nil {
			return nil, nil, errors.Wrap(err, "godotenv cannot load config")
		}
	}

	return ctx, cancel, nil
}
