package tx

// Transaction interface
type Trans interface {
	// Get the underlying transction instance
	GetPrincipal() interface{}
}

// Transaction callback function. return nil to commit the transaction, error to rollback the transaction
type TransFunc func(tx Trans) error

// Transaction manager, use it to start a transaction
type TransMgr interface {
	Transaction(f TransFunc) error
}
