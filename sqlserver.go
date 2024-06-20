package simutils

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectSQLServer(dbConn *DBConnection) (err error) {
	var (
		dsn                       = dbConn.DSN
		newLogger                 logger.Interface
		slowThreshold             time.Duration
		logLevel                  LogLevel
		colorful                  bool
		ignoreRecordNotFoundError bool
		parameterizedQueries      bool
	)

	if dsn == "" {
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			dbConn.User,
			dbConn.Pass,
			dbConn.Host,
			dbConn.Port,
			dbConn.DBName,
		)
	}

	if slowThreshold == 0 {
		slowThreshold = time.Second
	}
	if logLevel == 0 {
		logLevel = LogLevelInfo
	}

	if dbConn.Debug {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				// SlowThreshold slow SQL threshold
				SlowThreshold: slowThreshold,
				// LogLevel
				LogLevel: logger.LogLevel(logLevel),
				// Colorful
				Colorful: colorful,
				// IgnoreRecordNotFoundError ignore ErrRecordNotFound error for logger
				IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
				// ParameterizedQueries
				ParameterizedQueries: parameterizedQueries,
			},
		)
	}

	dbConn.DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   dbConn.SkipDefaultTransaction,
		FullSaveAssociations:                     dbConn.FullSaveAssociations,
		Logger:                                   newLogger,
		DryRun:                                   dbConn.DryRun,
		PrepareStmt:                              dbConn.PrepareStmt,
		DisableAutomaticPing:                     dbConn.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: dbConn.DisableForeignKeyConstraintWhenMigrating,
		IgnoreRelationshipsWhenMigrating:         dbConn.IgnoreRelationshipsWhenMigrating,
		DisableNestedTransaction:                 dbConn.DisableNestedTransaction,
		AllowGlobalUpdate:                        dbConn.AllowGlobalUpdate,
		QueryFields:                              dbConn.QueryFields,
		CreateBatchSize:                          dbConn.CreateBatchSize,
		TranslateError:                           dbConn.TranslateError,
	})

	return
}
