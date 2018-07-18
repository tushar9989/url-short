package database

import (
	"math/big"

	"github.com/tushar9989/url-short/write-server/internal/models"
)

type Database interface {
	Save(id big.Int, data models.LinkData) *DbError
	LoadServerMeta(name string) (models.ServerMeta, *DbError)
	UpdateServerCount(name string, count big.Int) *DbError
	Close()
}

type DbError struct {
	msg  string
	Code int
}

func (e *DbError) Error() string { return e.msg }
