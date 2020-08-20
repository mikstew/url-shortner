package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/catinello/base62"
	log "github.com/sirupsen/logrus"
)

const (
	shortUrl = "https://dol.ly/%s"
)

// Actual db would use auto-incrementing id as primary key
var db = make([]string, 0)

func shorten(w http.ResponseWriter, r *http.Request) {
	// Parse request for user input URL to shorten
	longUrl, err := parseShortenParam(r.URL.String())
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, err.Error())
		return
	}

	// Write user input URL to DB and get id in return
	id := writeToDb(longUrl)
	// Convert to base62 value that will be used for short URL
	token := base62.Encode(id)
	shortUrl := fmt.Sprintf(shortUrl, token)

	log.WithFields(log.Fields{
		"Long URL":  longUrl,
		"Short URL": shortUrl}).Info("URL Shortened")

	// Send response
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s => %s\n", longUrl, shortUrl)
}

func expand(w http.ResponseWriter, r *http.Request) {
	// Parse request for user token
	token, err := parsePathValue(r.URL.String())
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, err.Error())
		return
	}

	// Decode token to find id
	id, _ := base62.Decode(token)
	// Look for the long URL in db stored at id
	url, err := readFromDb(id)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, err.Error())
		return
	}

	log.WithFields(log.Fields{
		"Token": token,
		"URL":   url}).Info("URL Expanded")

	// Send response
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", url)
}

// writeToDb - "writes" to local slice acting as DB. Returns index of url.
func writeToDb(url string) int {
	db = append(db, url)
	return len(db) - 1
}

// readFromDb - "Reads" URL from DB stored at index id. Returns error if not found.
func readFromDb(id int) (string, error) {
	if id > len(db) {
		return "", fmt.Errorf("Value not found")
	}
	return db[id], nil
}

// parsePathValue - retrieves the token passed in as path
func parsePathValue(s string) (string, error) {
	// TODO: validate token as hex value
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("Error parsing token: %s", err)
	}
	// Return token without leading /
	return u.Path[1:], nil
}

// parseShortenParam - Parse the URL string and returns the url if provided as
// an input parameter
func parseShortenParam(s string) (string, error) {
	// TODO: sanitize user input to combat against SQL injection, XSS, etc.
	// URL encode & to ensure properly parsing user input URLs which contain &
	s = strings.ReplaceAll(s, "&", "%26")

	url, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("Error parsing 'url' parameter: %s", err)
	}

	q := url.Query()

	var longUrl string
	if len(q["url"]) > 0 {
		longUrl = q["url"][0]
	} else {
		return "", fmt.Errorf("'url' parameter is required")
	}

	return longUrl, nil
}

func handler() {
	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/", expand)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handler()
}
