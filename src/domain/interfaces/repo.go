package interfaces

// Repo generic interface
type Repo interface {
	FindOne(id string, e interface{}) error
	Upsert(id string, e interface{}) error
}
