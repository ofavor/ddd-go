package db

// Database interface
type Database interface {

	// Get connection due to the underlying implementation
	GetConn() interface{}

	// Register models
	RegisterModels([]interface{})
}
