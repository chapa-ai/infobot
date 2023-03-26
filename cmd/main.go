package main

import (
	"github.com/sirupsen/logrus"
	cfg "infoBot/internal/config"
	"infoBot/internal/db"
	"infoBot/internal/os"
	"infoBot/internal/telegram"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	logger := logrus.WithFields(logrus.Fields{})
	logger.Info("program started")

	err := cfg.InitConfig("/configs/config.env")
	if err != nil {
		logger.Errorf("init configs failed: %s", err)
		return err
	}

	connDb, err := db.InitDB()
	if err != nil {
		logger.Errorf("couldn't instantiate db: %s", err)
		return err
	}
	defer func() {
		err = connDb.Close()
		if err != nil {
			logger.Errorf("closing db failed: %s", err)
			return
		}
	}()

	err = db.MigrateDb("migrations")
	if err != nil {
		logger.Errorf("failed making migrations: %s", err)
		return err
	}
	logger.Info("migrations implemented")

	telegram.Run(logger)

	<-os.NotifyAboutExit()
	return nil
}
