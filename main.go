package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cespare/xxhash"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"sync"
)

func main() {
	// get websites
	urls, err := getSites("./sites.txt")
	if err != nil {
		fmt.Println("Could not read sites file: ", err)
		return
	}

	var wg sync.WaitGroup

	for _, u := range urls {
		// Increment the wait group counter
		wg.Add(1)
		go func(url string) {
			// Decrement the counter when the go routine completes
			defer wg.Done()
			// Call the function check
			hash, err := checkUrlHash(url)
			if err != nil {
				return
			}
			fmt.Println(hash)

		}(u)
	}
	// Wait for all the checkWebsite calls to finish
	wg.Wait()
}

// checks and prints a message if a website is up or down
func checkUrlHash(url string) (uint64, error) {
	// query the site
	resp, err := http.Get(url)
	fmt.Printf("Checking %s\n", url)
	if err != nil {
		fmt.Println(url, "is down !!!", err)
		return 0, err
	}

	// get the html content
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(url, "body could not be parsed:", err)
		return 0, err
	}

	// parse out the stuff we dont care about
	bodyString, err := extractBodyContent(body)
	if err != nil {
		fmt.Println(url, "body could not be parsed:", err)
		return 0, err
	}

	// hash body content
	hash := xxhash.Sum64String(bodyString)
	fmt.Println(hash)

	return hash, nil
}

// https://www.bacancytechnology.com/qanda/golang/extract-html-body-content-as-a-string-in-go
func extractBodyContent(htmlResponse []byte) (string, error) {
	respBody, err := html.Parse(bytes.NewReader(htmlResponse))
	if err != nil {
		return "", err
	}

	var bodyContent string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				var buf bytes.Buffer
				html.Render(&buf, c)
				bodyContent += buf.String()
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(respBody)

	return bodyContent, nil
}

func getSites(path string) ([]string, error) {
	sites := []string{}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening sites file:", err)
		return []string{}, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sites = append(sites, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading sites file with scanner:", err)
		return []string{}, err
	}

	return sites, nil
}
