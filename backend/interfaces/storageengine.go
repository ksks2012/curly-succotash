package interfaces

// StorageEngine defines methods of underlying storage
type StorageEngine interface {
	// Open setup and create underlying database connection pool and related resources.
	Open() (err error)

	// Close release underlying database connection pool and related resources.
	Close() (err error)
}
