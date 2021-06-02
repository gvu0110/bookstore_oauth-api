package cassandra

import (
	"github.com/gocql/gocql"
)

// QueryInterface should describe commonly used functions of the
// gocql.Query
type QueryInterface interface {
	Exec() error
	Iter() IterInterface
	Scan(...interface{}) error
}

// Query is a wrapper for a gocql.Query for mockability.
type Query struct {
	query *gocql.Query
}

// NewQuery instantiates a new Query
func NewQuery(query *gocql.Query) QueryInterface {
	return &Query{query}
}

// Exec wraps the query's Exec method
func (q *Query) Exec() error {
	return q.query.Exec()
}

// Iter wraps the query's Iter method
func (q *Query) Iter() IterInterface {
	return NewIter(q.query.Iter())
}

// Scan wraps the query's Scan method
func (q *Query) Scan(dest ...interface{}) error {
	return q.query.Scan(dest...)
}
