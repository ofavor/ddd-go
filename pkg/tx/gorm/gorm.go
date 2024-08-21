package tx

import (
	"github.com/ofavor/ddd-go/pkg/tx"

	"gorm.io/gorm"
)

// trans implementation based on gorm
type gormTrans struct {
	conn *gorm.DB
}

func (t *gormTrans) GetPrincipal() interface{} {
	return t.conn
}

type gormTransMgr struct {
	conn *gorm.DB
}

// Create Gorm based transaction manager
func NewTransMgr(db *gorm.DB) tx.TransMgr {
	return &gormTransMgr{conn: db}
}

// Start a transaction
func (tm *gormTransMgr) Transaction(f tx.TransFunc) error {
	if tm.conn == nil {
		panic("[trans-gorm] Failed to start transaction, no database connection")
	}
	dummy := func(tx *gorm.DB) error {
		return f(&gormTrans{conn: tx})
	}
	return tm.conn.Transaction(dummy)
}
