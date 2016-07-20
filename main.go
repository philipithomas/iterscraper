// iterscraper scrapes information from a website where URLs contain an incrementing integer.
// Information is retrieved from HTML5 elements, and outputted as a CSV.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var (
		urlTemplate = flag.String("url", "http://example.com/v/%d", "The URL you wish to scrape, containing \"%d\" where the id should be substituted")
		idLow       = flag.Int("from", 0, "The first ID that should be searched in the URL - inclusive.")
		idHigh      = flag.Int("to", 1, "The last ID that should be searched in the URL - exclusive")
		concurrency = flag.Int("concurrency", 1, "How many scrapers to run in parallel. (More scrapers are faster, but more prone to rate limiting or bandwith issues)")
		outfile     = flag.String("output", "output.csv", "Filename to export the CSV results")
		name        = flag.String("nameQuery", ".name", "JQuery-style query for the name element")
		address     = flag.String("addressQuery", ".address", "JQuery-style query for the address element")
		phone       = flag.String("phoneQuery", ".phone", "JQuery-style query for the phone element")
		email       = flag.String("emailQuery", ".email", "JQuery-style query for the email element")
	)
	flag.Parse()

	type task struct {
		url string
		id  int
	}
	tasks := make(chan task)
	go func() {
		for i := *idLow; i < *idHigh; i++ {
			tasks <- task{url: fmt.Sprintf(*urlTemplate, i), id: i}
		}
		close(tasks)
	}()

	sites := make(chan *site)
	var wg sync.WaitGroup
	wg.Add(*concurrency)
	go func() {
		wg.Wait()
		close(sites)
	}()

	for i := 0; i < *concurrency; i++ {
		go func() {
			defer wg.Done()
			for t := range tasks {
				site, err := fetch(t.url, t.id, *name, *address, *phone, *email)
				if err != nil {
					log.Printf("could not fetch %v: %v", t.url, err)
					continue
				}
				sites <- site
			}
		}()
	}

	dumpCSV(*outfile, sites)
}

func dumpCSV(path string, sites <-chan *site) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create file %s: %v", path, err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// Write headers to file
	if err := w.Write([]string{"id", "name", "url", "address", "phone", "email"}); err != nil {
		log.Fatalf("error writing record to csv: %v", err)
	}

	for s := range sites {
		if err := w.Write([]string{strconv.Itoa(s.id), s.name, s.url, s.address, s.phone, s.email}); err != nil {
			log.Fatalf("could not write record to csv: %v", err)
		}
	}

	if err := w.Error(); err != nil {
		return fmt.Errorf("writer failed: %v", err)
	}
	return nil
}

type site struct {
	url     string
	id      int
	address string
	phone   string
	email   string
	name    string
}

func fetch(url string, id int, nameQuery, addressQuery, phoneQuery, emailQuery string) (*site, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get %s: %v", url, err)
	}

	// GoQuery doesn't actually close the body - we have to do that
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Check for rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("you are being rate limited")
		}

		return nil, fmt.Errorf("bad response from server: %s", resp.Status)
	}

	// Load response into GoQuery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse page: %v", err)
	}

	// Pull info we want
	return &site{
		url:     url,
		id:      id,
		name:    strings.TrimSpace(doc.Find(nameQuery).Text()),
		address: strings.TrimSpace(doc.Find(addressQuery).Text()),
		phone:   strings.TrimSpace(doc.Find(phoneQuery).Text()),
		email:   strings.TrimSpace(doc.Find(emailQuery).Text()),
	}, nil
}
