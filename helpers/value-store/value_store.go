package valuestore

type ValueStore interface {
	Lookup(id string) ([]byte, error)
	Save(value string) (string, error)
	// possibly add a clear?
	// possibly add a total num entries?
}
