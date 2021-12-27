package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// ConnectDB estabelece a conexão com o banco de dados.
func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbMaxIddleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDDLE_CONNS"))
	dbMaxOpensConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPENS_CONNS"))

	dns := fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=%v password=%v",
		dbHost, dbPort, dbUser, dbName, dbSSLMode, dbPassword)

	database, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
			NameReplacer:  strings.NewReplacer("CID", "Cid"),
		},
	})
	if err != nil {
		log.Fatalf("error to connect to database: %v", err)
	}

	db = database

	config, err := db.DB()
	if err != nil {
		log.Fatal("error:", err.Error())
	}

	config.SetMaxIdleConns(dbMaxIddleConns)
	config.SetMaxOpenConns(dbMaxOpensConns)
	config.SetConnMaxLifetime(time.Hour)
}

// GetDB retorna o banco de dados.
func GetDB() *gorm.DB {
	return db
}

// CloseDB encerra a conexão com o banco de dados.
func CloseDB() {
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatalf("failed to close connection to database: %v", err)
	}

	err = dbSQL.Close()
	if err != nil {
		log.Fatalf("failed to close connection to database: %v", err)
	}
}
