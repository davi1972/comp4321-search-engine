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

## Quick Start

Run the program:
```bash
$ go run main.go
```

Print out result
```bash
$ go run test.go
```

## Specification
Written in Go Programming Language using databse BadgerDB

## People
This program is made by David Sun and Hans Krishandi


