package simutils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	// errors
	ErrInvalidDatabaseConnection = errors.New("invalid database connection")
	ErrInvalidDatabaseDriver     = errors.New("invalid database driver used")
	ErrConnectionAlreadyExist    = errors.New("database connection already exist")

	// vars
	defaultDBName = "default"
	pool          []*DBConnection
)

type DBs map[string]*DBConnection

type DBConnection struct {
	DBConfig `mapstructure:",squash"`
	// Deprecated: use DSN instead
	Name string `json:"name,omitempty" mapstructure:"name"`
	// Deprecated: use DSN instead
	Host string `json:"host,omitempty" mapstructure:"host"`
	// Deprecated: use DSN instead
	Port string `json:"port,omitempty" mapstructure:"port"`
	// Deprecated: use DSN instead
	User string `json:"user,omitempty" mapstructure:"user"`
	// Deprecated: use DSN instead
	Pass string `json:"pass,omitempty" mapstructure:"pass"`
	// Deprecated: use DSN instead
	DBName string `json:"db_name,omitempty" mapstructure:"db_name"`

	DB *gorm.DB `json:"-"`
}

type DBConfig struct {
	// Driver is one of
	// 1: PostgresSQL, 2: SQLServer, 3: SQLite
	Driver DatabaseDriver `json:"driver,omitempty" mapstructure:"driver"`
	// DSN is database connection string
	DSN string `json:"dsn,omitempty" mapstructure:"dsn"`
	// Debug enables debugger
	Debug bool `json:"debug,omitempty" mapstructure:"debug"`
	// GORM perform single create, update, delete operations in transactions by default to ensure database data integrity
	// You can disable it by setting `SkipDefaultTransaction` to true
	SkipDefaultTransaction bool `json:"skip_default_transaction,omitempty" mapstructure:"skip_default_transaction"`
	// FullSaveAssociations full save associations
	FullSaveAssociations bool `json:"full_save_associations,omitempty" mapstructure:"full_save_associations"`
	// Logger
	Logger LoggerConfig `json:"logger,omitempty" mapstructure:"logger"`
	// DryRun generate sql without execute
	DryRun bool `json:"dry_run,omitempty" mapstructure:"dry_run"`
	// PrepareStmt executes the given query in cached statement
	PrepareStmt bool `json:"prepare_stmt,omitempty" mapstructure:"prepare_stmt"`
	// DisableAutomaticPing
	DisableAutomaticPing bool `json:"disable_automatic_ping,omitempty" mapstructure:"disable_automatic_ping"`
	// DisableForeignKeyConstraintWhenMigrating
	DisableForeignKeyConstraintWhenMigrating bool `json:"disable_foreign_key_constraint_when_migrating,omitempty" mapstructure:"disable_foreign_key_constraint_when_migrating"`
	// IgnoreRelationshipsWhenMigrating
	IgnoreRelationshipsWhenMigrating bool `json:"ignore_relationships_when_migrating,omitempty" mapstructure:"ignore_relationships_when_migrating"`
	// DisableNestedTransaction disable nested transaction
	DisableNestedTransaction bool `json:"disable_nested_transaction,omitempty" mapstructure:"disable_nested_transaction"`
	// AllowGlobalUpdate allow global update
	AllowGlobalUpdate bool `json:"allow_global_update,omitempty" mapstructure:"allow_global_update"`
	// QueryFields executes the SQL query with all fields of the table
	QueryFields bool `json:"query_fields,omitempty" mapstructure:"query_fields"`
	// CreateBatchSize default create batch size
	CreateBatchSize int `json:"create_batch_size,omitempty" mapstructure:"create_batch_size"`
	// TranslateError enabling error translation
	TranslateError bool `json:"translate_error,omitempty" mapstructure:"translate_error"`
}

// LoggerConfig logger config
type LoggerConfig struct {
	SlowThreshold             time.Duration `json:"slow_threshold,omitempty" mapstructure:"slow_threshold"`
	Colorful                  bool          `json:"colorful,omitempty" mapstructure:"colorful"`
	IgnoreRecordNotFoundError bool          `json:"ignore_record_not_found_error,omitempty" mapstructure:"ignore_record_not_found_error"`
	ParameterizedQueries      bool          `json:"parameterized_queries,omitempty" mapstructure:"parameterized_queries"`
	LogLevel                  LogLevel      `json:"log_level,omitempty" mapstructure:"log_level"`
}

// LogLevel log level
type LogLevel int

const (
	// Silent silent log level
	LogLevelSilent LogLevel = iota + 1
	// Error error log level
	LogLevelError
	// Warn warn log level
	LogLevelWarn
	// Info info log level
	LogLevelInfo
)

func (c DBConnection) Value() (value driver.Value, err error) {
	var b []byte
	if b, err = json.Marshal(c); err != nil {
		return
	}

	return string(b), nil
}

func (c *DBConnection) Scan(value interface{}) (err error) {
	if value == nil {
		*c = DBConnection{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := DBConnection{}
	if err = json.Unmarshal(bytes, &result); err != nil {
		return
	}

	*c = result

	return nil
}

func (c *DBConnection) IsValid() bool {
	return c.Driver != 0
}

func Add(dbConn *DBConnection) (err error) {
	if !dbConn.IsValid() {
		return ErrInvalidDatabaseConnection
	}

	for _, c := range pool {
		if dbConn.Name == c.Name {
			return ErrConnectionAlreadyExist
		}
	}

	if dbConn.DB == nil {
		if err = Connect(dbConn); err != nil {
			return
		}
	}

	pool = append(pool, dbConn)

	return nil
}

func Get(connecitonName string) (dbConn *DBConnection, err error) {
	for _, c := range pool {
		if connecitonName == c.Name {
			return c, nil
		}
	}

	return nil, ErrNotFound
}

func ConnectDBs(dbs map[string]*DBConnection) error {
	if len(dbs) == 0 {
		return nil
	}

	var (
		defaultDB = dbs[defaultDBName]
	)

	if defaultDB != nil {
		var dbresolvers = &dbresolver.DBResolver{}

		logrus.Infof("connecting to %s", defaultDBName)
		if err := Connect(defaultDB); err != nil {
			return err
		}

		for dbname, db := range dbs {
			if dbname == defaultDBName {
				continue
			}

			relatedTables := viper.GetStringSlice(fmt.Sprintf("databases.%s.related_to", dbname))

			rels := make([]any, len(relatedTables))
			for idx, rt := range relatedTables {
				rels[idx] = rt
			}

			dbresolvers = dbresolvers.Register(
				dbresolver.Config{
					Sources: []gorm.Dialector{dialector(db)},
				}, rels...,
			)
		}

		if err := defaultDB.DB.Use(dbresolvers); err != nil {
			return err
		}
	}

	for dbname, db := range dbs {
		if dbname == defaultDBName {
			continue
		}

		logrus.Infof("connecting to %s", dbname)
		if err := Connect(db); err != nil {
			return err
		}
	}

	return nil
}

func Connect(dbConn *DBConnection) (err error) {
	var (
		d                         gorm.Dialector
		newLogger                 logger.Interface
		slowThreshold             time.Duration
		logLevel                  LogLevel
		colorful                  bool
		ignoreRecordNotFoundError bool
		parameterizedQueries      bool
	)

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

	if d = dialector(dbConn); d == nil {
		return ErrInvalidDatabaseDriver
	}

	dbConn.DB, err = gorm.Open(d, &gorm.Config{
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

	return err
}

func dialector(dbConn *DBConnection) gorm.Dialector {
	var (
		dsn = dbConn.DSN
	)

	switch dbConn.Driver {
	case PostgresSQL:
		if dsn == "" {
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				dbConn.Host,
				dbConn.User,
				dbConn.Pass,
				dbConn.DBName,
				dbConn.Port,
			)
		}
		return postgres.Open(dsn)
	case SQLServer:
		if dsn == "" {
			dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
				dbConn.User,
				dbConn.Pass,
				dbConn.Host,
				dbConn.Port,
				dbConn.DBName,
			)
		}
		return sqlserver.Open(dsn)
	case SQLite:
		// Connect to the database
		// sqliteDB, err := sql.Open("sqlite3", dsn)
		// if err != nil {
		// 	panic("failed to connect database")
		// }
		// defer sqliteDB.Close()

		// // Add the REGEXP function
		// addRegexpFunction(sqliteDB)
		// return sqlite.Dialector{Conn: sqliteDB}
		return sqlite.Open(dsn)
	}

	return nil
}
