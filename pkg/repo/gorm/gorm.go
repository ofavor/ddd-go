package gorm

import (
	"github.com/ofavor/ddd-go/pkg/entity"
	"github.com/ofavor/ddd-go/pkg/repo"
	"github.com/ofavor/ddd-go/pkg/tx"

	"gorm.io/gorm"
)

// gorm repository
type GormRepo[E entity.Entity[D], D any] struct {
	conn   *gorm.DB
	loader EntityLoader[E, D]
}

type EntityLoader[E entity.Entity[D], D any] func(d *D) E

// Create gorm repository
func NewRepo[E entity.Entity[D], D any](conn *gorm.DB, loader EntityLoader[E, D]) *GormRepo[E, D] {
	return &GormRepo[E, D]{
		conn:   conn,
		loader: loader,
	}
}

// Get real connection
func (r *GormRepo[E, D]) GetConn(tx tx.Trans) *gorm.DB {
	if tx == nil {
		return r.conn
	}
	if conn, ok := tx.GetPrincipal().(*gorm.DB); !ok {
		panic("[repo-gorm] Transaction principal is not a *gorm.DB instance")
	} else {
		return conn
	}
}

func (r *GormRepo[E, D]) prepareQuery(query *gorm.DB, filter repo.Filter, sorter repo.Sorter) {
	if filter != nil {
		conds := filter.Conditions()
		for k, v := range conds {
			switch vt := v.(type) {
			case []interface{}:
				query = query.Where(k, vt...)
			default:
				query = query.Where(k, v)
			}
		}
	}
	if sorter != nil {
		sorts := sorter.Sorts()
		for _, v := range sorts {
			query = query.Order(v)
		}
	}
}

// Count implements repo.Repository.
func (r *GormRepo[E, D]) Count(tx tx.Trans, filter repo.Filter) (cnt int64, err error) {
	conn := r.GetConn(tx)
	query := conn.Model(new(D))
	r.prepareQuery(query, filter, nil)
	err = query.Count(&cnt).Error
	return
}

// List implements repo.Repository.
func (r *GormRepo[E, D]) List(tx tx.Trans, filter repo.Filter, sorter repo.Sorter, offset int64, limit int64) ([]E, error) {
	conn := r.GetConn(tx)
	query := conn.Model(new(D))
	r.prepareQuery(query, filter, sorter)
	arr := make([]*D, 0)
	if err := query.Offset(int(offset)).Limit(int(limit)).Find(&arr).Error; err != nil {
		return nil, err
	}
	out := make([]E, 0, len(arr))
	for _, a := range arr {
		out = append(out, r.loader(a))
	}
	return out, nil
}

// Get implements repo.Repository.
func (r *GormRepo[E, D]) Get(tx tx.Trans, id interface{}) (e E, err error) {
	conn := r.GetConn(tx)
	m := new(D)
	if err = conn.First(m, id).Error; err != nil {
		return
	}
	e = r.loader(m)
	return
}

// Save implements repo.Repository.
func (r *GormRepo[E, D]) Save(tx tx.Trans, e E) error {
	conn := r.GetConn(tx).Session(&gorm.Session{FullSaveAssociations: true})
	if e.IsNew() {
		return conn.Create(e.DAO()).Error
	} else {
		return conn.Save(e.DAO()).Error
	}
}

// Delete implements repo.Repository.
func (r *GormRepo[E, D]) Delete(tx tx.Trans, id interface{}) error {
	conn := r.GetConn(tx)
	return conn.Delete(new(D), id).Error
}
