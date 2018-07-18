package database

import (
	"math/big"

	"github.com/tushar9989/url-short/read-server/internal/models"
)

type Database interface {
	LoadLinkData(id big.Int) (models.LinkData, error)
	LoadLinkStatsForUser(userId string) []models.LinkStats
	IncrementLinkStatsForUser(userId string, linkId big.Int) error
	Close()
}
