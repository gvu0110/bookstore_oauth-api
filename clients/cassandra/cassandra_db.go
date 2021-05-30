package cassandra

import (
	"os"

	"github.com/gocql/gocql"
)

const (
	cassandra_address_env_var = "CASSANDRA_ADDRESS"
)

var (
	session *gocql.Session
)

func init() {
	// Connect to the Cassandra cluster
	cluster := gocql.NewCluster(os.Getenv(cassandra_address_env_var))
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
}

func GetSession() *gocql.Session {
	return session
}
