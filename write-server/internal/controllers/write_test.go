package controllers

import (
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestUrl(t *testing.T) {
	body := strings.NewReader(`{
			"TargetUrl": "abx",
			"ExpireAt": "2006-01-02T15:04:05Z"
		}`)
	_, err := getLinkDataFromBody(body, "")

	if err == nil {
		t.Error("Passed with invalid URL")
	}

	body = strings.NewReader(`{
			"TargetUrl": "http://google.com",
			"ExpireAt": "2006-01-02T15:04:05Z"
		}`)
	_, err = getLinkDataFromBody(body, "")

	if err != nil {
		t.Error("Failed with valid URL: ", err)
	}
}

func TestInvalidPayload(t *testing.T) {
	body := strings.NewReader(`sdbj`)
	_, err := getLinkDataFromBody(body, "")

	if err == nil {
		t.Error("Passed with invalid body")
	}
}

func TestUserId(t *testing.T) {
	body := strings.NewReader(`{
		"TargetUrl": "http://google.com",
		"ExpireAt": "2006-01-02T15:04:05Z"
	}`)

	linkData, _ := getLinkDataFromBody(body, "")

	if linkData.UserId != "-1" {
		t.Error(`Expected "-1", got `, linkData.UserId)
	}

	body = strings.NewReader(`{
		"TargetUrl": "http://google.com",
		"ExpireAt": "2006-01-02T15:04:05Z"
	}`)
	linkData, _ = getLinkDataFromBody(body, "abc")

	if linkData.UserId != "abc" {
		t.Error(`Expected "abc", got `, linkData.UserId)
	}
}

func TestEmailsWithoutUser(t *testing.T) {
	body := strings.NewReader(`{
		"TargetUrl": "http://google.com",
		"ExpireAt": "2006-01-02T15:04:05Z",
		"ValidEmails": [
			"tushar9989@gmail.com"
		]
	}`)

	_, err := getLinkDataFromBody(body, "")

	if err == nil {
		t.Error("Saved emails without valid user id")
	}
}

func TestExpireTime(t *testing.T) {
	body := strings.NewReader(`{
		"TargetUrl": "http://google.com",
		"ExpireAt": "2006-01-02T15:04:05Z"
	}`)

	linkData, _ := getLinkDataFromBody(body, "")
	expected, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	if linkData.ExpireAt.String() != expected.String() {
		t.Error("Expected ", expected.String(), "got ", linkData.ExpireAt.String())
	}

	body = strings.NewReader(`{
		"TargetUrl": "http://google.com"
	}`)

	_, err := getLinkDataFromBody(body, "")

	if err == nil {
		t.Error("Passed without ExpireAt")
	}
}

func TestEncodeDecodeId(t *testing.T) {
	number := big.NewInt(1675108645995)

	slug := encodeIdToSlug(*number)

	if slug != "tushar1" {
		t.Error("Expected tushar1, got ", slug)
	}

	decoded, _ := decodeSlugToId(slug)

	if decoded.Cmp(number) != 0 {
		t.Error("Expected 1675108645995, got ", decoded.String())
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
