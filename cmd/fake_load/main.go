package main

import (
	"bufio"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/starnuik/golang_messagio/internal"
	"github.com/starnuik/golang_messagio/internal/cmd"
)

var words []string
var client *resty.Client
var endpoint string

func randBetween(min int, max int) int {
	return min + int(rand.Int31n(int32(max-min)))
}

func readWords(path string) []string {
	file, err := os.Open(path)
	cmd.ServerError(err)
	defer file.Close()

	words := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) < 3 {
			continue
		}
		words = append(words, word)
	}

	// shrink, just in case
	words = append([]string{}, words...)
	return words
}

func randomWords(from []string, count int) string {
	words := []string{}
	for range count {
		idx := rand.Int31n(int32(len(from)))
		words = append(words, from[idx])
	}
	return strings.Join(words, " ")
}

func postRequest() {
	wordCount := randBetween(3, 25)
	payload := randomWords(words, wordCount)
	req := internal.MessageRequest{Content: payload}

	res, err := client.R().
		SetBody(req).
		Post(endpoint)
	cmd.ServerError(err)

	if res.StatusCode() != http.StatusOK {
		log.Println(endpoint, "endpoint failed with", res.StatusCode(), ", body:", string(res.Body()))
	}
}

func main() {
	endpoint = os.Getenv("SERVICE_MESSAGE_URL")
	words = readWords("./words.txt")
	client = resty.New()

	// send rand() requests every 60 seconds, somewhat evenly spaced
	for {
		messageDensity := randBetween(100, 1000)
		delay := time.Millisecond * time.Duration(60*1000/messageDensity)

		log.Println("starting request batch: density", messageDensity, ", delay", delay)
		for range messageDensity {
			go postRequest()
			time.Sleep(delay)
		}
	}
}
