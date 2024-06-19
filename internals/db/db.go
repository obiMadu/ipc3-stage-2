package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/obimadu/ipc3-stage-2/internals/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var counts int
var DB *gorm.DB

func InitDB() {
	// new db
	db := connectToMysql()

	if gin.Mode() == gin.ReleaseMode {
		db.Logger.LogMode(0)
	}

	DB = db

	rawDB := RawDB()

	rawDB.SetMaxIdleConns(20)
	rawDB.SetMaxOpenConns(100)

	// migrate models
	err := migrate()
	if err != nil {
		log.Panicf("Unable to migrate models %s\n", err.Error())
	}
	log.Println("Successfully Migrated Models.")
}

func migrate() error {
	err := DB.AutoMigrate(&models.Users{})
	if err != nil {
		return err
	}

	return nil
}

func RawDB() *sql.DB {
	rawDB, err := DB.DB()
	if err != nil {
		log.Panicf("Unable to get raw sql.DB %s\n", err.Error())
	}

	return rawDB
}

func connectToPostgres() *gorm.DB {
	dsn := os.Getenv("POSTGRES_DSN")

	for {
		connection, err := openPostgres(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Fatal(err)
		}

		log.Println("Backing off for three seconds....")
		time.Sleep(3 * time.Second)
		continue
	}
}

func connectToMysql() *gorm.DB {
	dsn := os.Getenv("MYSQL_DSN")

	for {
		connection, err := openMysql(dsn)
		if err != nil {
			log.Println("MySQL not yet ready ...")
			counts++
		} else {
			log.Println("Connected to MySQL")
			return connection
		}

		if counts > 10 {
			log.Fatal(err)
		}

		log.Println("Backing off for three seconds....")
		time.Sleep(3 * time.Second)
		continue
	}
}

func openPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// return *sql.DB from db(*gorm.DB) to enable Ping()
	gormDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// ping database
	err = gormDB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func openMysql(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// return *sql.DB from db(*gorm.DB) to enable Ping()
	gormDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// ping database
	err = gormDB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
