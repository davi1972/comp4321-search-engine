# COMP4321 Search Engine Project

This program is made to complete COMP 4321 Project requirements to make a working search engine consisted of a spider to fetch pages recursively and indexer to extract keywords from a page.

## Installation

Install Go

```bash
$ wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz
$ sudo tar -C /usr/local/ -xzf go1.10.linux-amd64.tar.gz
$ echo "export PATH=\$PATH:/usr/local/go/bin" | sudo tee -a /etc/profile
$ source /etc/profile
```

Download the repository
```bash
$ go get github.com/davi1972/comp4321-search-engine
```

Download all the dependencies

```bash
$ go get -u github.com/gocolly/colly/...
$ go get -u github.com/reiver/go-porterstemmer/...
$ go get github.com/dgraph-io/badger/...
```

Navigate to the frontend folder. Change the API_URL to the backend inside the script.js file.

## Quick Start

Run the indexer:
```bash
$ go run indexer.go
```

To change the depth and starting page, edit the *maxDepth* and *rootPage* variable in the indexer.go

Run the backend
```bash
$ go run backend.go
```

Will run backend serving in http://localhost:8000

## Specification
Written in Go Programming Language using database BadgerDB

## People
This program is made by David Sun, Hans Krishandi, and Calvin Cheng
