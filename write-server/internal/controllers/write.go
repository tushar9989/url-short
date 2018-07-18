package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tushar9989/url-short/write-server/internal/database"
	"github.com/tushar9989/url-short/write-server/internal/models"

	"github.com/tushar9989/url-short/write-server/internal/pkg/counters"
)

func Write(counter counters.BigInt, db database.Database) func(http.ResponseWriter, *http.Request) {
	return wrapper(func(r *http.Request) (interface{}, error) {
		if r.Method != "POST" {
			return nil, errors.New("Invalid request method")
		}

		linkData, err := getLinkDataFromBody(r.Body, r.Header.Get("X-User-ID"))

		if err != nil {
			return nil, err
		}

		diff := linkData.ExpireAt.Sub(time.Now().UTC())
		if diff.Minutes() <= 5 || diff.Hours() > 30*24 {
			return nil, errors.New("Invalid expire time")
		}

		if linkData.CustomSlug != "" {
			id, err := decodeSlugToId(linkData.CustomSlug)
			if err != nil {
				return nil, err
			}
			dbErr := db.Save(id, linkData)
			if dbErr == nil {
				return map[string]string{"slug": linkData.CustomSlug}, nil
			}
			return nil, dbErr
		} else {
			id := counter.IncrementAndGetOldValue()
			dbErr := db.Save(id, linkData)
			for dbErr != nil && dbErr.Code == 3 {
				id = counter.IncrementAndGetOldValue()
				dbErr = db.Save(id, linkData)
			}
			if dbErr == nil {
				return map[string]string{"slug": encodeIdToSlug(id)}, nil
			}
			return nil, dbErr
		}
	})
}

func encodeIdToSlug(input big.Int) string {
	number := big.NewInt(0)
	number.Add(number, &input)
	base := big.NewInt(62)
	zero := big.NewInt(0)
	arr := make([]string, 0)

	for number.Cmp(zero) != 0 {
		mod := big.NewInt(0)
		mod.Mod(number, base)

		intMod := mod.Int64()

		if intMod < 10 {
			arr = append(arr, string('0'+intMod))
		} else if intMod < 36 {
			intMod -= 10
			arr = append(arr, string('a'+intMod))
		} else {
			intMod -= 36
			arr = append(arr, string('A'+intMod))
		}

		number.Div(number, base)
	}

	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}

	return strings.Join(arr, "")
}

func decodeSlugToId(number string) (big.Int, error) {
	base := big.NewInt(62)
	multiplier := big.NewInt(1)
	answer := big.NewInt(0)

	if len(number) > 7 {
		return *big.NewInt(0), errors.New("Input slug too long")
	}

	for i := len(number) - 1; i >= 0; i-- {
		var intVal *big.Int
		if number[i] >= '0' && number[i] <= '9' {
			intVal = big.NewInt(int64(number[i] - '0'))
		} else if number[i] >= 'a' && number[i] <= 'z' {
			intVal = big.NewInt(int64(10 + number[i] - 'a'))
		} else if number[i] >= 'A' && number[i] <= 'Z' {
			intVal = big.NewInt(int64(36 + number[i] - 'A'))
		} else {
			return *big.NewInt(0), errors.New("Invalid character in slug")
		}

		answer.Add(answer, intVal.Mul(intVal, multiplier))

		multiplier.Mul(multiplier, base)
	}

	return *answer, nil
}

func wrapper(h func(r *http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(map[string]string{"status": "FAIL", "message": fmt.Sprintf("%v", r)})
			}
		}()
		response, err := h(r)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"status": "FAIL", "message": err.Error()})
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK", "data": response})
		}
	}
}

func getLinkDataFromBody(reader io.Reader, userId string) (models.LinkData, error) {
	decoder := json.NewDecoder(reader)
	var linkData models.LinkData
	err := decoder.Decode(&linkData)

	if err != nil {
		return linkData, err
	}

	if linkData.ExpireAt == nil {
		return linkData, errors.New("ExpireAt must be set")
	}

	if userId == "" {
		userId = "-1"
		if len(linkData.ValidEmails) > 0 {
			return linkData, errors.New("Valid Emails cannot be set unless logged in.")
		}
	}
	linkData.UserId = userId

	_, err = url.ParseRequestURI(linkData.TargetUrl)
	if err != nil {
		return linkData, err
	}

	return linkData, nil
}
