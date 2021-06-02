package cassandra

import (
	"github.com/gocql/gocql"
)

// IterInterface should describe commonly used functions of the
// gocql.Iter
type IterInterface interface {
	Close() error
	Scan(...interface{}) bool
}

// Iter is a wrapper for an gocql.Iter for mockability.
type Iter struct {
	iter *gocql.Iter
}

// NewIter instantiates a new Iter
func NewIter(iter *gocql.Iter) IterInterface {
	return &Iter{iter}
}

// Close is a wrapper for the iter's Close method
func (i *Iter) Close() error {
	return i.iter.Close()
}

// Scan is a wrapper for the iter's Scan method
func (i *Iter) Scan(dest ...interface{}) bool {
	return i.iter.Scan(dest...)
}
