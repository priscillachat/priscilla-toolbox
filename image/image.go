package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	imageAPI = "https://www.googleapis.com/customsearch/v1"
	gifAPI   = "http://api.giphy.com/v1/gifs/search"
)

type imgRespBody struct {
	Items []imgLink `json:"items"`
}

type imgLink struct {
	Link string `json:"link"`
}

type gifRespBody struct {
	Data []gifImgs `json:"data"`
}

type gifImgs struct {
	Images gifData `json:"images"`
}

type gifData struct {
	Original map[string]string `json:"original"`
}

func main() {
	mode := flag.String("mode", "image", "Mode: image or gif")
	safe := flag.String("safe", "medium",
		"(image mode only) Safe search: high/medium/off")
	cseID := flag.String("cseid", "", "Google custom search ID")
	cseKey := flag.String("csekey", "", "Google custom search Key")
	giphyKey := flag.String("giphykey", "dc6zaTOxFJmzC",
		"Giphy api key (leave default to use public beta key)")

	query := ""
	flag.StringVar(&query, "query", "", "search query")

	flag.Parse()

	if *mode != "image" && *mode != "gif" {
		fmt.Println("Mode must be either 'image' or 'gif'")
		os.Exit(1)
	}

	if *mode == "image" && (*safe != "high" && *safe != "medium" && *safe != "off") {
		fmt.Println("safe must be either 'high', 'medium', or 'off'")
		os.Exit(1)
	}

	if *mode == "image" && (*cseID == "" || *cseKey == "") {
		fmt.Println("Missing cseid or csekey for image mode")
		os.Exit(1)
	}

	if query == "" {
		fmt.Println("You didn't say what to search, so grumpy cat it is.")
		query = "grumpycat"
	}

	htclient := &http.Client{Timeout: time.Duration(10) * time.Second}
	rand.Seed(time.Now().Unix())

	if *mode == "image" {
		req, err := http.NewRequest("GET", imageAPI, nil)
		q := req.URL.Query()
		queryParams := map[string]string{
			"v":          "1.0",
			"searchType": "image",
			"q":          query,
			"safe":       *safe,
			"fields":     "items(link)",
			"rsz":        "8",
			"cx":         *cseID,
			"key":        *cseKey,
		}
		for key, val := range queryParams {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
		resp, err := htclient.Do(req)

		if err != nil {
			fmt.Println("Query error:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body := imgRespBody{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&body)
			if err != nil {
				fmt.Println("Decoding error:", err)
				os.Exit(1)
			}
			resultCount := len(body.Items)
			if resultCount < 1 {
				fmt.Println("That sounds gibberish, try a different search")
				os.Exit(1)
			}
			fmt.Println(body.Items[rand.Intn(resultCount)].Link)
		} else {
			fmt.Println("Couldn't get image from google")
			os.Exit(1)
		}
	} else {
		req, err := http.NewRequest("GET", gifAPI, nil)
		q := req.URL.Query()
		q.Add("q", query)
		q.Add("api_key", *giphyKey)
		req.URL.RawQuery = q.Encode()
		resp, err := htclient.Do(req)

		if err != nil {
			fmt.Println("Query error:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body := gifRespBody{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&body)
			if err != nil {
				fmt.Println("Decoding error:", err)
				os.Exit(1)
			}
			resultCount := len(body.Data)
			if resultCount < 1 {
				fmt.Println("That sounds gibberish, try a different search")
				os.Exit(1)
			}
			fmt.Println(body.Data[rand.Intn(resultCount)].Images.Original["url"])
		} else {
			fmt.Println("Problem talking to giphy")
			os.Exit(1)
		}

	}
}
