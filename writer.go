package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
)

// Write fetched data to CSV
func writeSites(dataChan chan site, wg *sync.WaitGroup) {
	file, err := os.Create(outfile)
	if err != nil {
		panic(fmt.Sprintf("Unable to open file %s - error %s", outfile, err))
	}

	defer file.Close()

	w := csv.NewWriter(file)

	// Write headers to file
	if err := w.Write(dataRowHeaders); err != nil {
		log.Fatalf("error writing record to csv: %v", err)
	}

	for {
		s := <-dataChan
		if s.err != nil {
			// There was a problem with this site
			log.Printf("Unable to fetch id %d - %s", s.id, s.err)
		} else {
			// There was not a problem with this site
			log.Printf("Fetched id %d", s.id)
			if err := w.Write(s.dataRow()); err != nil {
				log.Fatalf("error writing record to csv: %v", err)
			}
			w.Flush()
			if err := w.Error(); err != nil {
				log.Fatal(err)
			}

		}
		// Decrement waitgroup so that we know when to stop writing to the output
		wg.Done()
	}
}
