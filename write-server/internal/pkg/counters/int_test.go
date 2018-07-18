package counters_test

import (
	"math/big"
	"testing"

	"github.com/tushar9989/url-short/write-server/internal/pkg/counters"
)

func TestInvalidConfig(t *testing.T) {
	_, err := counters.NewBigInt(*big.NewInt(11), *big.NewInt(20), *big.NewInt(10))

	if err == nil {
		t.Error("Current less than start failed")
	}

	_, err = counters.NewBigInt(*big.NewInt(11), *big.NewInt(0), *big.NewInt(10))

	if err == nil {
		t.Error("Current greater than end failed")
	}
}

func TestIncrement(t *testing.T) {
	counter, _ := counters.NewBigInt(*big.NewInt(10), *big.NewInt(20), *big.NewInt(10))

	_ = counter.IncrementAndGetOldValue()
	value := counter.IncrementAndGetOldValue()

	if value.Cmp(big.NewInt(11)) != 0 {
		t.Error("Expected 11, got ", value.String())
	}
}

func TestValue(t *testing.T) {
	counter, _ := counters.NewBigInt(*big.NewInt(10), *big.NewInt(20), *big.NewInt(10))

	value := counter.Value()
	_ = counter.IncrementAndGetOldValue()

	if value.Cmp(big.NewInt(10)) != 0 {
		t.Error("Expected 10, got ", value.String())
	}
}

func TestRollOver(t *testing.T) {
	counter, _ := counters.NewBigInt(*big.NewInt(10), *big.NewInt(20), *big.NewInt(20))

	_ = counter.IncrementAndGetOldValue()
	_ = counter.IncrementAndGetOldValue()
	value := counter.IncrementAndGetOldValue()

	if value.Cmp(big.NewInt(11)) != 0 {
		t.Error("Expected 11, got ", value.String())
	}
}

func BenchmarkIncrement(b *testing.B) {
	counter, _ := counters.NewBigInt(*big.NewInt(0), *big.NewInt(100), *big.NewInt(20))
	for n := 0; n < b.N; n++ {
		_ = counter.IncrementAndGetOldValue()
	}
}
