# iterscraper

[![Build Status](https://travis-ci.org/philipithomas/iterscraper.svg?branch=master)](https://travis-ci.org/philipithomas/iterscraper)

A basic package used for scraping information from a website where URLs contain an incrementing integer. Information is retrieved from HTML5 elements, and outputted as a CSV.

## Flags

Flags are all optional, and are set with a single dash on the command line, e.g.

```
iterscraper \
-url            "http://foo.com/%d" \
-from           1                   \
-to             10                  \
-concurrency    10                  \
-output         foo.csv             \
-nameQuery      ".name"             \
-addressQuery   ".address"          \
-phoneQuery     ".phone"            \
-emailQuery     ".email"            
```

For an explanation of the options, type `iterscraper -help`

General usage of iterscraper:

```
  -addressQuery string
        JQuery-style query for the address element (default ".address")
  -concurrency int
        How many scrapers to run in parallel. (More scrapers are faster, but more prone to rate limiting or bandwith issues) (default 1)
  -emailQuery string
        JQuery-style query for the email element (default ".email")
  -from int
        The first ID that should be searched in the URL - inclusive.
  -nameQuery string
        JQuery-style query for the name element (default ".name")
  -output string
        Filename to export the CSV results (default "output.csv")
  -phoneQuery string
        JQuery-style query for the phone element (default ".phone")
  -to int
        The last ID that should be searched in the URL - exclusive (default 1)
  -url string
        The URL you wish to scrape, containing "%d" where the id should be substituted (default "http://example.com/v/%d")
```

## URL Structure

Successive pages must look like:

```
http://example.com/foo/1/bar
http://example.com/foo/2/bar
http://example.com/foo/3/bar
```

iterscraper would then accept the url in the following style, in `Printf` style such that numbers may be substituted into the url:

```
http://example.com/foo/%d/bar
```

## Installation

Building the source requires the [Go programming language](https://golang.org/doc/install) and the [Glide](http://glide.sh) package manager.

```
# Dependency is GoQuery
go get github.com/PuerkitoBio/goquery
# Get and build source
go get github.com/philipithomas/iterscraper

# If your $PATH is configured correctly, you can call it directly
iterscraper [flags]

```


## Errata

* This is purpose-built for some internal scraping. It's not meant to be the scraping tool for every user case, but you're welcome to modify it for your purposes
* On a `429 - too many requests` error, the app stops. This is because data integrity (being able to tell whether or not a page exists) is compromised. One possible fix is to turn down concurrency.
* The package will [follow up to 10 redirects](https://golang.org/pkg/net/http/#Get)
* On a `404 - not found` error, the system will log the miss, then continue. It is not exported to the CSV.

