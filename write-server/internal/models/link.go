package models

import (
	"math/big"
	"time"
)

type LinkData struct {
	Id          big.Int
	ExpireAt    *time.Time
	TargetUrl   string
	ValidEmails []string
	UserId      string
	CustomSlug  string
}
