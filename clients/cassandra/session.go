package cassandra

import (
	"os"

	"github.com/gocql/gocql"
)

const (
	cassandra_address_env_var = "CASSANDRA_ADDRESS"
)

// SessionInterface should describe commonly used functions of the
// gocql.Session
type SessionInterface interface {
	Close()
	Query(string, ...interface{}) QueryInterface
}

// Sessions is a wrapper for a docql.Session for mockability
type Session struct {
	session *gocql.Session
}

// NewSession instantiates a new Session
func NewSession(session *gocql.Session) SessionInterface {
	return &Session{session: session}
}

// Close wraps the session's close method
func (s *Session) Close() {
	s.session.Close()
}

// Query wraps the session's query method
func (s *Session) Query(stmt string, values ...interface{}) QueryInterface {
	return NewQuery(s.session.Query(stmt, values...))
}

func NewConnection() SessionInterface {
	cluster := gocql.NewCluster(os.Getenv(cassandra_address_env_var))
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	return NewSession(session)
}
