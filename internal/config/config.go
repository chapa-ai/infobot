package configs

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func InitConfig(path string) error {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Errorf("couldn't get mainPath: %s", err)
		return err
	}
	err = godotenv.Load(filepath.Join(dir, path))
	if err != nil {
		logrus.Errorf("failed godotenv.Load: %s", err)
		return err
	}

	return nil
}
