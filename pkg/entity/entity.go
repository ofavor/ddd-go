package entity

type PersistSupport[D any] interface {
	IsNew() bool
	DAO() *D
}

type Entity[D any] interface {
}
