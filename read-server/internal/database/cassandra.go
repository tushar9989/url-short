package database

import (
	"math/big"

	"github.com/gocql/gocql"
	"github.com/tushar9989/url-short/read-server/internal/models"
)

type Cassandra struct {
	session *gocql.Session
}

func NewCassandra(servers []string, keyspace string) (*Cassandra, error) {
	cluster := gocql.NewCluster(servers...)
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &Cassandra{session: session}, nil
}

func (c *Cassandra) Close() {
	if !c.session.Closed() {
		c.session.Close()
	}
}

func (c *Cassandra) LoadLinkData(id big.Int) (models.LinkData, error) {
	returnValue := models.LinkData{}

	if err := c.session.Query(`SELECT id, user_id, allowed_emails, expire_time, target_url FROM links WHERE id = ? LIMIT 1`,
		id.String()).Consistency(gocql.One).Scan(&returnValue.Id, &returnValue.UserId, &returnValue.ValidEmails, &returnValue.ExpireAt, &returnValue.TargetUrl); err != nil {
		return returnValue, err
	}

	return returnValue, nil
}

func (c *Cassandra) LoadLinkStatsForUser(userId string) []models.LinkStats {
	returnValue := make([]models.LinkStats, 0)

	iter := c.session.Query("SELECT id, views FROM views WHERE user_id = ?", userId).Iter()
	var Id big.Int
	var Views int64
	for iter.Scan(&Id, &Views) {
		IdCopy := *big.NewInt(0).Set(&Id)
		returnValue = append(returnValue, models.LinkStats{Id: &IdCopy, Views: Views})
	}

	return returnValue
}

func (c *Cassandra) IncrementLinkStatsForUser(userId string, linkId big.Int) error {
	err := c.session.Query(`UPDATE views
		SET views = views + 1
		WHERE id = ? 
		AND user_id = ?;`,
		linkId.String(),
		userId,
	).Consistency(gocql.Quorum).Exec()

	return err
}
