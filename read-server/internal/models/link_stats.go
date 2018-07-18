package models

import (
	"math/big"
)

type LinkStats struct {
	Id    *big.Int
	Views int64
	Slug  string
}
