package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DavidHuie/gomigrate"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"infoBot/internal/models"
	"os"
	"sync"
	"time"
)

var (
	db         *sql.DB
	dbConnOnce sync.Once
)

func InitDB() (*sql.DB, error) {

	var err error

	dbConnOnce.Do(func() {
		host := os.Getenv("POSTGRES_HOST")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbName := os.Getenv("POSTGRES_DB")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)
		d, _err := sql.Open("postgres", dsn)

		if _err != nil {
			logrus.Errorf("sql.Open failed: %s\n", _err)
			err = _err
			return
		}
		db = d

	})

	return db, err
}

func MigrateDb(path string) error {
	db, err := InitDB()
	if err != nil {
		logrus.Errorf("failed GetDB(): %s", err)
		return err
	}
	migrator, err := gomigrate.NewMigrator(db, gomigrate.Postgres{}, path)
	if err != nil {
		logrus.Errorf("failed implement migrations: %s", err)
		return err
	}

	return migrator.Migrate()
}

func SaveCurrencies(ctx context.Context, currencyBTCUSDT *models.Data) (*models.ResponseData, error) {
	currency := &models.ResponseData{}

	err := db.QueryRowContext(ctx, `INSERT INTO currencies("symbol", "buy", "time") VALUES($1, $2, $3) RETURNING "symbol", "buy", "time"`,
		currencyBTCUSDT.Symbol, currencyBTCUSDT.Buy, time.Now()).Scan(&currency.Symbol, &currency.Buy, &currency.Time)
	if err != nil {
		logrus.Errorf("failed db.QueryRowContext: %s", err)
		return nil, err
	}

	return currency, nil
}

func TimeOfFirstQuery(ctx context.Context, id int) (string, error) {
	var t string
	err := db.QueryRowContext(ctx, `SELECT time FROM currencies WHERE "id" = $1`, id).Scan(&t)
	if err != nil {
		logrus.Errorf("failed db.QueryRowContext: %s", err)
	}

	return t, nil
}

func CountOfAllQueries(ctx context.Context) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM currencies`).Scan(&count)
	if err != nil {
		logrus.Errorf("failed db.QueryRowContext: %s", err)
	}

	return count, nil
}
