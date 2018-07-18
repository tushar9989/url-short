package models

import (
	"math/big"
)

type ServerMeta struct {
	Name    string
	Start   big.Int
	End     big.Int
	Current big.Int
}
