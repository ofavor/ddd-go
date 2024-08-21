package gorm

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/ofavor/ddd-go/pkg/db"
	"github.com/ofavor/ddd-go/pkg/util"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gormDatabase database implementation based on gorm
type gormDatabase struct {
	conn *gorm.DB
}

// Create gorm database
func NewDatabase(
	driver string,
	dns string,
	encKey string,
	debug bool,
) db.Database {
	l := logger.Warn
	if debug {
		l = logger.Info
	}
	if strings.Trim(encKey, " ") != "" {
		encryptionKey = encKey
	}
	var conn *gorm.DB
	var err error
	conf := &gorm.Config{
		Logger:                                   logger.Default.LogMode(l),
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	switch driver {
	case "mysql":
		conn, err = gorm.Open(mysql.Open(dns), conf)
	case "postgres":
		conn, err = gorm.Open(postgres.Open(dns), conf)
	default:
		panic("Unsupported database driver: " + driver)
	}
	if err != nil {
		panic(err)
	}
	return NewDatabaseWithConn(conn)
}

func NewDatabaseWithConn(conn *gorm.DB) db.Database {
	return &gormDatabase{conn}
}

// Get connection returns *gorm.DB
func (d *gormDatabase) GetConn() interface{} {
	return d.conn
}

// Register models, gorm will generate tables automatically
func (d *gormDatabase) RegisterModels(models []interface{}) {
	if d.conn == nil {
		panic("No database connection")
	}
	err := d.conn.
		// Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(models...)
	if err != nil {
		panic(err)
	}
}

// Encrypted table column key, must not be changed if used
var encryptionKey = "S20jBJE0b71GdKnP"

// Encrypted table column
type Encrypted string

// Scan implement gorm interface
func (e *Encrypted) Scan(value interface{}) error {
	h := ""
	switch v := value.(type) {
	case []byte:
		h = string(v)
	case string:
		h = v
	default:
		return fmt.Errorf("value must be string: %s", value)
	}
	str, err := util.AesDecrypt([]byte(encryptionKey), h)
	if err != nil {
		return err
	}
	*e = Encrypted(str)
	return nil
}

// Scan implement gorm interface
func (e Encrypted) Value() (driver.Value, error) {
	str, err := util.AesEncrypt([]byte(encryptionKey), string(e))
	if err != nil {
		return nil, err
	}
	return str, nil
}
