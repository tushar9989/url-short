package database

import (
	"math/big"
	"time"

	"github.com/gocql/gocql"
	"github.com/tushar9989/url-short/write-server/internal/models"
)

type Cassandra struct {
	session *gocql.Session
}

func NewCassandra(servers []string, keyspace string) (*Cassandra, *DbError) {
	cluster := gocql.NewCluster(servers...)
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, &DbError{msg: err.Error(), Code: 1}
	}

	return &Cassandra{session: session}, nil
}

func (c *Cassandra) Close() {
	if !c.session.Closed() {
		c.session.Close()
	}
}

func (c *Cassandra) Save(id big.Int, data models.LinkData) *DbError {
	data.Id = id
	diff := data.ExpireAt.Sub(time.Now().UTC())
	if diff.Seconds() <= 0 || diff.Hours() > 30*24 {
		return &DbError{msg: "Invalid expire time", Code: 2}
	}
	seconds := int64(diff.Seconds())
	existingData := make(map[string]interface{})
	applied, err := c.session.Query(`INSERT INTO 
		links (id, expire_time, target_url, allowed_emails, user_id) 
		VALUES (?, ?, ?, ?, ?)
		IF NOT EXISTS
		USING TTL ?`,
		id,
		data.ExpireAt,
		data.TargetUrl,
		data.ValidEmails,
		data.UserId,
		seconds,
	).Consistency(gocql.Quorum).MapScanCAS(existingData)

	if !applied {
		return &DbError{msg: "Already exists", Code: 3}
	}

	if err == nil {
		return nil
	}

	return &DbError{msg: err.Error(), Code: 1}
}

func (c *Cassandra) LoadServerMeta(name string) (models.ServerMeta, *DbError) {
	returnValue := models.ServerMeta{}

	if err := c.session.Query(`SELECT name, current, start, end FROM write_server_meta WHERE name = ? LIMIT 1`,
		name).Consistency(gocql.One).Scan(&returnValue.Name, &returnValue.Current, &returnValue.Start, &returnValue.End); err != nil {
		return returnValue, &DbError{msg: err.Error(), Code: 1}
	}

	return returnValue, nil
}

func (c *Cassandra) UpdateServerCount(name string, count big.Int) *DbError {
	existingData := make(map[string]interface{})
	applied, err := c.session.Query(`UPDATE 
		write_server_meta 
		SET current = ? 
		WHERE name = ? 
		IF EXISTS`,
		count,
		name,
	).Consistency(gocql.Quorum).MapScanCAS(existingData)

	if !applied {
		return &DbError{msg: "Details for given name not found", Code: 1}
	}

	if err == nil {
		return nil
	}

	return &DbError{msg: err.Error(), Code: 1}
}
