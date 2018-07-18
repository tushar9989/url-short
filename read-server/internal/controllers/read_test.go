package controllers

import (
	"math/big"
	"testing"
)

func TestDecodeSlug(t *testing.T) {
	number := big.NewInt(5621662037)

	decoded, _ := decodeSlugToId("68rV6l")

	if decoded.Cmp(number) != 0 {
		t.Error("Expected 5621662037, got ", decoded.String())
	}
}

func TestInvalidSlug(t *testing.T) {

	slug := "_*()"

	_, err := decodeSlugToId(slug)

	if err == nil {
		t.Error("Passed with invalid slug")
	}

	slug = "ZZZZZZZZ"

	_, err = decodeSlugToId(slug)

	if err == nil {
		t.Error("Passed with long slug")
	}
}
