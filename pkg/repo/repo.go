package repo

import (
	"github.com/ofavor/ddd-go/pkg/entity"
	"github.com/ofavor/ddd-go/pkg/tx"
)

// Filter interface, provide query conditions
type Filter interface {
	Conditions() map[string]interface{}
}

// Sorter interface, provide sort conditions
type Sorter interface {
	Sorts() []string
}

// Repository interface
type Repository[E entity.Entity[D], D any] interface {
	// Count number of records
	Count(tx tx.Trans, filter Filter) (int64, error)
	// List records
	List(tx tx.Trans, filter Filter, sorter Sorter, offset, limit int64) ([]E, error)
	// Get a record by id
	Get(tx tx.Trans, id interface{}) (E, error)
	// Save a record
	Save(tx tx.Trans, e E) error
	// Delete a record
	Delete(tx tx.Trans, id interface{}) error
}
