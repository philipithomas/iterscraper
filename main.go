package main

import (
	"flag"
	"sync"
)

var (
	urlBase       string
	idLow, idHigh int
	concurrency   int
	outfile       string
	nameQuery     string
	addressQuery  string
	phoneQuery    string
	emailQuery    string
)

func main() {
	// Get flags
	flag.StringVar(&urlBase, "url", "http://example.com/v/%d", `The URL you wish to scrape, containing "%d" where the id should be substituted`)
	flag.IntVar(&idLow, "from", 0, "The first ID that should be searched in the URL - inclusive.")
	flag.IntVar(&idHigh, "to", 1, "The last ID that should be searched in the URL - exclusive")
	flag.IntVar(&concurrency, "concurrency", 1, "How many scrapers to run in parallel. (More scrapers are faster, but more prone to rate limiting or bandwith issues)")
	flag.StringVar(&outfile, "output", "output.csv", "Filename to export the CSV results")
	flag.StringVar(&nameQuery, "nameQuery", ".name", "JQuery-style query for the name element")
	flag.StringVar(&addressQuery, "addressQuery", ".address", "JQuery-style query for the address element")
	flag.StringVar(&phoneQuery, "phoneQuery", ".phone", "JQuery-style query for the phone element")
	flag.StringVar(&emailQuery, "emailQuery", ".email", "JQuery-style query for the email element")

	flag.Parse()

	// Use waitgroup so we can keep track of tasks
	var wg sync.WaitGroup
	wg.Add(idHigh - idLow)

	// channel for emitting sites to fetch
	taskChan := make(chan site)
	// Channel of data to write to disk
	dataChan := make(chan site)

	go emitTasks(taskChan)

	for i := 0; i < concurrency; i++ {
		go scrape(taskChan, dataChan)
	}

	go writeSites(dataChan, &wg)

	wg.Wait()
}
