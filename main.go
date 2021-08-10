package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/harishdurga/go-pinger/db"

	"go.mongodb.org/mongo-driver/bson"
)

type Result struct {
	url        string
	err        error
	latency    time.Duration
	statusCode int
}

func ping(url string, ch chan<- Result) {
	start := time.Now()
	if resp, err := http.Get(url); err != nil {
		ch <- Result{url, err, 0, resp.StatusCode}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		ch <- Result{url, nil, t, resp.StatusCode}
		resp.Body.Close()
	}
}

func main() {
	urls := []string{}
	//Opening the CSV file containing the URLs
	csvData, err := os.Open("ping_urls.csv")
	if err != nil {
		panic(err)
	}
	defer csvData.Close()

	//Creating a new CSV reader
	r, err := csv.NewReader(csvData).ReadAll()
	if err != nil {
		panic(err)
	}
	//Iterating through the CSV file
	for i, row := range r {
		if i == 0 {
			continue
		}
		urls = append(urls, row[0])
	}

	resultsCh := make(chan Result)
	//Creating a new goroutine for each URL
	for _, url := range urls {
		go ping(url, resultsCh)
	}

	//Results to be uploaded to mongodb
	var documents []interface{}
	excutedTime := time.Now()

	//Waiting for all goroutines to finish
	for range urls {
		result := <-resultsCh
		documents = append(documents, bson.M{
			"url":         result.url,
			"err":         result.err,
			"latency":     (result.latency),
			"executed":    excutedTime.String(),
			"status_code": result.statusCode,
		})
	}
	client, ctx, cancel, err := db.Connect()
	if err != nil {
		panic(err)
	}
	insertManyResult, err := db.InsertMany(client, ctx, "pinger",
		"pings", documents)

	// handle the error
	if err != nil {
		panic(err)
	}

	fmt.Println("Result of InsertMany")

	// print the insertion ids of the multiple
	// documents, if they are inserted.
	for id := range insertManyResult.InsertedIDs {
		fmt.Println(id)
	}
	// Release resource when main function is returned.
	defer db.Close(client, ctx, cancel)
}
