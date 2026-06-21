package database

import (
	"auth-server/config"
	"auth-server/model"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connect to db with connection pooling, retry, and health check
func ConnectDB() {
	dsn := buildDSN()

	var err error
	for i := 1; i <= 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/5): %v", i, err)
		time.Sleep(time.Duration(i) * 2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect database after 5 attempts: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	configurePool(sqlDB)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connection Opened to Database")
	model.RunAutoMigrate(DB)
}

func buildDSN() string {
	host := config.Config("DB_HOST")
	portStr := config.Config("DB_PORT")
	user := config.Config("DB_USER")
	password := config.Config("DB_PASSWORD")
	dbname := config.Config("DB_NAME")
	sslmode := config.Config("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	port, err := strconv.ParseUint(portStr, 10, 32)
	if err != nil {
		log.Fatalf("Failed to parse DB_PORT value: %v", err)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

func configurePool(sqlDB *sql.DB) {
	maxOpenStr := config.Config("DB_MAX_OPEN_CONNS")
	maxIdleStr := config.Config("DB_MAX_IDLE_CONNS")
	maxLifetimeStr := config.Config("DB_CONN_MAX_LIFETIME")
	maxIdleTimeStr := config.Config("DB_CONN_MAX_IDLETIME")

	if v, err := strconv.Atoi(maxOpenStr); err == nil && v > 0 {
		sqlDB.SetMaxOpenConns(v)
	} else {
		sqlDB.SetMaxOpenConns(25)
	}

	if v, err := strconv.Atoi(maxIdleStr); err == nil && v >= 0 {
		sqlDB.SetMaxIdleConns(v)
	} else {
		sqlDB.SetMaxIdleConns(10)
	}

	if d, err := time.ParseDuration(maxLifetimeStr); err == nil && d > 0 {
		sqlDB.SetConnMaxLifetime(d)
	} else {
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	if d, err := time.ParseDuration(maxIdleTimeStr); err == nil && d > 0 {
		sqlDB.SetConnMaxIdleTime(d)
	} else {
		sqlDB.SetConnMaxIdleTime(1 * time.Minute)
	}
}
