package cassandra

import (
	"github.com/gocql/gocql"
)

var (
	cluster *gocql.ClusterConfig
)

func init() {
	// Connect to the Cassandra cluster
	cluster = gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

}

func GetSession() (*gocql.Session, error) {
	return cluster.CreateSession()
}
