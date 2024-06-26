package simutils

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectPostgres(dbConn *DBConnection) (err error) {
	var newLogger logger.Interface
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConn.Host,
		dbConn.User,
		dbConn.Pass,
		dbConn.DBName,
		dbConn.Port,
	)

	if dbConn.Debug {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
				// IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			},
		)
	}

	dbConn.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	return
}
