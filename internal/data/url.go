package data

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DOWNLOAD_URL = `https://en-word.net/static/english-wordnet-$year.xml.gz`

// finds newest DB file to download
func FindMostRecentUrl() (string, error) {
	year := time.Now().Year()

	for i := range 10 {
		currentUrl := strings.Replace(DOWNLOAD_URL, "$year", strconv.Itoa(year-i), 1)
		response, err := http.Get(currentUrl)
		if err != nil {
			return "", err
		}

		if response.StatusCode != http.StatusOK {
			continue
		}

		return currentUrl, nil
	}

	return "", fmt.Errorf("could not find any version of database")
}
