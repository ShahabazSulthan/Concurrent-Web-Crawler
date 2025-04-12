package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	MaxConcurrency = 5                  
	RequestTimeout = 5 * time.Second    
	URLFile        = "urls.txt"         
)

type Result struct {
	URL    string
	Status int
	Title  string
	Error  error
}

func main() {
	start := time.Now()

	urls, err := readURLs(URLFile)
	if err != nil {
		log.Fatalf("couldn't read url file: %v", err)
	}

	results := crawlURLs(urls)

	fmt.Println("AUthor: Shahabaz Sultha")

	for _, res := range results {
		if res.Error != nil {
			log.Printf("[ERROR] %s: %v", res.URL, res.Error)
		} else {
			log.Printf("[OK] %s (%d) -> %s", res.URL, res.Status, res.Title)
		}
	}

	log.Printf("✅ Done in %v", time.Since(start))
	fmt.Println("Thank You")
}

func readURLs(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	urlMap := make(map[string]struct{})
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		url := strings.TrimSpace(line)
		if url != "" {
			urlMap[url] = struct{}{}
		}
	}

	unique := make([]string, 0, len(urlMap))
	for url := range urlMap {
		unique = append(unique, url)
	}
	return unique, nil
}

func crawlURLs(urls []string) []Result {
	var wg sync.WaitGroup
	resultChan := make(chan Result, len(urls))
	semaphore := make(chan struct{}, MaxConcurrency)

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			semaphore <- struct{}{}        
			defer func() { <-semaphore }() 

			resultChan <- fetchURLWithTimeout(u)
		}(url)
	}

	wg.Wait()
	close(resultChan)

	results := make([]Result, 0, len(urls))
	for res := range resultChan {
		results = append(results, res)
	}
	return results
}

func fetchURLWithTimeout(url string) Result {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Error: err}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Result{URL: url, Error: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK || resp.Body == nil {
		return Result{
			URL:    url,
			Status: resp.StatusCode,
			Error:  fmt.Errorf("got status code %d", resp.StatusCode),
		}
	}

	title, err := extractTitle(resp.Body)
	if err != nil {
		return Result{URL: url, Status: resp.StatusCode, Error: err}
	}

	return Result{
		URL:    url,
		Status: resp.StatusCode,
		Title:  title,
	}
}

func extractTitle(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	var title string
	var findTitle func(*html.Node)

	findTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil && title == ""; c = c.NextSibling {
			findTitle(c)
		}
	}
	findTitle(doc)

	if title == "" {
		return "No Title Found", nil
	}
	return title, nil
}
