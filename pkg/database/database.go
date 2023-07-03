package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"strconv"
)

type Database struct {
	Conn *sqlx.DB
}

func InitDB(bbddName string) *Database {
	var err error
	db := &Database{}

	db.Conn, err = SqlLiteConnect(bbddName)

	if err != nil {
		log.Error("Error connecting database: ", err)
	} else {
		log.Info("Database connected")
	}

	return db
}

func SqlLiteConnect(bbddName string) (*sqlx.DB, error) {
	dir, _ := os.Getwd()
	db, err := sqlx.Connect("sqlite", filepath.Dir(dir)+bbddName)

	numbSc, err := GetDBVersion(db)
	if numbSc < len(scripts)-1 {
		err = CreateScripts(db, numbSc)
		if err != nil {
			return db, err
		}
	}
	log.Info("Database version :", len(scripts))

	db.SetMaxOpenConns(3)

	return db, err
}

func CreateScripts(db *sqlx.DB, numbSc int) error {

	for i := numbSc; i < len(scripts); i++ {
		log.Info("Executing script - ", i+1)
		_, err := db.Exec(scripts[i].Script)
		if err != nil {
			log.Error(err)
			return err
		}
		err = UpdateVersion(db, i)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	log.Info("Scrips executed successfully - ", len(scripts))

	return nil
}

func UpdateVersion(db *sqlx.DB, version int) error {
	_, err := db.Exec("UPDATE db_version SET version = " + strconv.Itoa(version) + " WHERE 1 = 1;")
	return err
}

func GetDBVersion(db *sqlx.DB) (int, error) {
	var versions []int
	err := db.Select(&versions, "SELECT * FROM db_version")
	if err != nil {
		return 0, err
	} else {
		return versions[0], err
	}
}

func RemoveDB(bbddName string) error {
	dir, _ := os.Getwd()
	return os.Remove(filepath.Dir(dir) + bbddName)
}
