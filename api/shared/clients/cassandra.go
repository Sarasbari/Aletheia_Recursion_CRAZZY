package clients

import (
	"context"
	"time"

	"github.com/gocql/gocql"
)

type CassandraLogger struct {
	session *gocql.Session
}

func NewCassandraLogger(hosts []string, keyspace string) (*CassandraLogger, error) {
	if len(hosts) == 0 || keyspace == "" {
		return &CassandraLogger{}, nil
	}
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &CassandraLogger{session: session}, nil
}

func (c *CassandraLogger) LogEvent(_ context.Context, eventType, entityID, payload string) error {
	if c.session == nil {
		return nil
	}
	return c.session.Query(
		"INSERT INTO request_logs (event_type, entity_id, payload, created_at) VALUES (?, ?, ?, ?)",
		eventType, entityID, payload, time.Now().UTC(),
	).Exec()
}

func (c *CassandraLogger) Close() {
	if c.session != nil {
		c.session.Close()
	}
}
