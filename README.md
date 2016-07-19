# iterscraper

A basic package used for scraping information from a website where URLs contain an integer. Information is retrieved from HTML5 elements, and outputted as a CSV.

## Flags


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
# Get and build source
go get github.com/philipithomas/iterscraper
# Install dependencies
glide install
# If your $PATH is configured correctly, you can call it directly
iterscraper [flags]

```


## Errata

* This is purpose-built for some internal scraping. It's not meant to be the scraping tool for every user case, but you're welcome to modify it for your purposes
* On a `429 - too many requests` error, the app panics. This is because data integrity (being able to tell whether or not a page exists) is compromised. One possible fix is to turn down concurrency.
* The package will [follow up to 10 redirects](https://golang.org/pkg/net/http/#Get)
* On a `404 - not found` error, the system will log the miss, then continue. It is not exported to the CSV.

