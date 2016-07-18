package main

import (
	"flag"
	"sync"
)

func main() {
	// Get flags
	urlBase := flag.String("url", "http://example.com/v/%d", `The URL you wish to scrape, containing "%d" where the id should be substituted`)
	iterLow := flag.Int("from", 0, "The first ID that should be searched in the URL - inclusive.")
	iterHigh := flag.Int("to", 1, "The last ID that should be searched in the URL - exclusive")
	concurrency := flag.Int("concurrency", 1, "How many scrapers to run in parallel. (More scrapers are faster, but more prone to rate limiting or bandwith issues)")
	outfile := flag.String("output", "output.csv", "Filename to export the CSV results")

	// make some fetchers
	var wg sync.WaitGroup

    // channel for emitting IDs
	idChan := make(chan int)

    // channel for confirming failed scrapes - to keep waitgroup synced
    errChan := make(chan bool)

    scrapedChan := make(chan )
}
