package main

import (
	"fmt"
	"github.com/tybeller/webcheck/web"
	"sync"
)

func main() {
	// get websites
	urls, err := web.GetSites("../sites.txt")
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
			hash, err := web.CheckUrlHash(url)
			if err != nil {
				return
			}
			fmt.Println(hash)

		}(u)
	}
	// Wait for all the checkWebsite calls to finish
	wg.Wait()
}
