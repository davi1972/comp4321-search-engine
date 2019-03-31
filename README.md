# COMP4321 Search Engine Project

This program is made to complete COMP 4321 Project requirements to make a working search engine consisted of a spider to fetch pages recursively and indexer to extract keywords from a page.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Quick Start

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
```

Running the program
```bash
$ go run main.go
```

Printing out result
```bash
$ go run test.go
```
## Specification
Written in Go Programming Language

## People
This program is made by David Sun and Hans Krishandi


