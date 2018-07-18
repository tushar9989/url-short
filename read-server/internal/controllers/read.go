package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/tushar9989/url-short/read-server/internal/database"
)

type ApiError struct {
	Message string
	Code    int
}

func ReadSlug(db database.Database) func(http.ResponseWriter, *http.Request) {
	return wrapper(func(w http.ResponseWriter, r *http.Request) (interface{}, *ApiError) {
		if r.Method != "GET" {
			return nil, &ApiError{"Invalid request method", http.StatusMethodNotAllowed}
		}

		userEmail := r.Header.Get("X-User-Email")

		slug := r.URL.Path[3:]

		id, err := decodeSlugToId(slug)

		if err != nil {
			return nil, &ApiError{err.Error(), http.StatusBadRequest}
		}

		linkData, err := db.LoadLinkData(id)

		if err != nil || linkData.ExpireAt.Before(time.Now().UTC()) {
			return nil, &ApiError{err.Error(), http.StatusNotFound}
		}

		if len(linkData.ValidEmails) != 0 && !contains(linkData.ValidEmails, userEmail) {
			return nil, &ApiError{"Not authorized.", http.StatusForbidden}
		}

		if linkData.UserId != "-1" {
			go db.IncrementLinkStatsForUser(linkData.UserId, id)
		}

		http.Redirect(w, r, linkData.TargetUrl, 301)

		return nil, nil
	})
}

func contains(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

func wrapper(h func(w http.ResponseWriter, r *http.Request) (interface{}, *ApiError)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(map[string]string{"status": "FAIL", "message": fmt.Sprintf("%v", r)})
			}
		}()
		response, err := h(w, r)
		if err != nil {
			w.WriteHeader(err.Code)
			json.NewEncoder(w).Encode(map[string]string{"status": "FAIL", "message": err.Message})
		} else {
			if response != nil {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK", "data": response})
			}
		}
	}
}
