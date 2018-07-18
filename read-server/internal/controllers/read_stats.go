package controllers

import (
	"math/big"
	"net/http"
	"strings"

	"github.com/tushar9989/url-short/read-server/internal/database"
)

func ReadUserStats(db database.Database) func(http.ResponseWriter, *http.Request) {
	return wrapper(func(w http.ResponseWriter, r *http.Request) (interface{}, *ApiError) {
		if r.Method != "GET" {
			return nil, &ApiError{"Invalid request method", http.StatusMethodNotAllowed}
		}

		userId := r.Header.Get("X-User-ID")

		if userId == "-1" || userId == "" {
			return nil, &ApiError{"Must be logged in to look at stats.", http.StatusBadRequest}
		}

		slugStats := db.LoadLinkStatsForUser(userId)

		for i, slugStat := range slugStats {
			slugStat.Slug = encodeIdToSlug(*slugStat.Id)
			slugStat.Id = nil
			slugStats[i] = slugStat
		}

		return slugStats, nil
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
