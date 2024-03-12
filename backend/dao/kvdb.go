package dao

import (
	"gobot/backend/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
)

var (
	MainDB *bolt.DB
	LogDB  *bolt.DB
)

func InitKvDB() error {
	path := viper.GetString("kvdb.path")
	if !utils.PathExist(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	mainDB, err := bolt.Open(filepath.Join(path, "main.db"), 0666, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return err
	}
	MainDB = mainDB

	logDB, err := bolt.Open(filepath.Join(path, "log.db"), 0666, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return err
	}
	LogDB = logDB

	return nil
}
