# parallel-http

> DRAFT VERSION

Simple tool firing parallel requests.

Intended to be used for API concurrency handling tests.

## Requirements

* [Go](https://golang.org/dl/) to be installed

## Build

    go build

## Usage

    ./parallel-http --help

## Example

    ./parallel-http -u http://localhost/v1/ -m PATCH -p 10 -h 'Content-Type:application/json\nX-Origin:test' -b "$(cat content.json)"
