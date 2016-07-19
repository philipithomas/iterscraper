package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type site struct {
	url     string
	body    []byte
	id      int
	address string
	phone   string
	email   string
	name    string
	err     error
}

// Fetch the website body
func (s *site) fetch() {
	resp, err := http.Get(s.url)
	if err != nil {
		s.err = err
		return
	}

	// GoQuery doesn't actually close the body - we have to do that
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Check for rate limiting
		if resp.StatusCode == 429 {
			log.Fatal("You are being rate limited - program exiting")
		}

		s.err = fmt.Errorf("Bad response from server - status code %d\n", resp.StatusCode)
		return
	}

	// Load response into GoQuery
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	// Pull info we want
	s.name = strings.TrimSpace(doc.Find(nameQuery).Text())
	s.address = strings.TrimSpace(doc.Find(addressQuery).Text())
	s.phone = strings.TrimSpace(doc.Find(phoneQuery).Text())
	s.email = strings.TrimSpace(doc.Find(phoneQuery).Text())

	// Clean up address data because it often has line breaks
	//s.address = new:ineRe.ReplaceAllString(s.address, " replace")
}

// Include headers with the row output format so that we can compare easily.
var dataRowHeaders = []string{"id", "name", "url", "address", "phone", "email"}

func (s *site) dataRow() []string {
	return []string{fmt.Sprintf("%d", s.id), s.name, s.url, s.address, s.phone, s.email}
}
