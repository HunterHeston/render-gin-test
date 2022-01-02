package valuestore

import "context"

type ValueStore interface {
	LookUp(ctx context.Context, id string) ([]byte, error)
	Save(ctx context.Context, value []byte) (string, error)
	// possibly add a clear?
	// possibly add a total num entries?
}
