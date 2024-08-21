package entity

type Entity[D any] interface {
	IsNew() bool
	DAO() *D
}
