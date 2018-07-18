package counters

import (
	"errors"
	"math/big"
	"sync"
)

type bigInt struct {
	mu    sync.Mutex
	value big.Int
	end   big.Int
	start big.Int
}

type BigInt interface {
	IncrementAndGetOldValue() (oldValue big.Int)
	Value() (value big.Int)
}

var one = big.NewInt(1)

func NewBigInt(start big.Int, end big.Int, value big.Int) (BigInt, error) {

	if value.Cmp(&start) < 0 || value.Cmp(&end) > 0 {
		return &bigInt{}, errors.New("Invalid counter configuration")
	}

	return &bigInt{value: value, end: end, start: start}, nil
}

func (c *bigInt) IncrementAndGetOldValue() (oldValue big.Int) {
	c.mu.Lock()
	oldValue = *big.NewInt(0).Set(&c.value)
	c.value.Add(&c.value, one)
	if c.value.Cmp(&c.end) > 0 {
		c.value = c.start
	}
	c.mu.Unlock()
	return
}

func (c *bigInt) Value() (value big.Int) {
	c.mu.Lock()
	value = *big.NewInt(0).Set(&c.value)
	c.mu.Unlock()
	return
}
