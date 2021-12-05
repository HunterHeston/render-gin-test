package valuestore

type ValueStore interface {
	LookUp(id string) ([]byte, error)
	Save(value []byte) (string, error)
	// possibly add a clear?
	// possibly add a total num entries?
}
