# Systembolaget beer releases

[![Build Status](https://travis-ci.org/jonasf/systembolaget-beer-releases.svg?branch=master)](https://github.com/jonasf/systembolaget-beer-releases)

This is a simple web site that use the data indexed by https://github.com/jonasf/systembolaget-article-indexer and displays which beers that will be released on Systembolaget.

## Getting started

Simplest way is use [Docker](https://www.docker.com/) and [Docker-compose](https://github.com/docker/compose)

1. Download the docker-compose.yml
2. Run `docker-compose up`

## Building

1. Make sure [Go](https://golang.org/) is installed and set up properly
2. Install [dep](https://github.com/golang/dep)
3. Clone the repo
4. Run `make all`
