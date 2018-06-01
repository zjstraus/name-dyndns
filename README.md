# name-dyndns [![Build Status](https://travis-ci.org/mfycheng/name-dyndns.svg?branch=master)](https://travis-ci.org/mfycheng/name-dyndns) [![GoDoc](https://godoc.org/github.com/mfycheng/name-dyndns?status.svg)](https://godoc.org/github.com/mfycheng/name-dyndns)
Client that automatically updates name.com DNS records.

## Getting name-dyndns

Since name-dyndns has no external dependencies, you can get it simply by:

```go
go get github.com/mfycheng/name-dyndns
```

## Requirements

In order to use name-dyndns, you must have an API key from name.com, which
can be requested from https://www.name.com/reseller/apply.

Once you have your API key, all you must do is setup `config.json`. An example
`config.json` file can be found in `api/config_test.json`.

## Running

name-dyndns will run in an infinite loop, constantly making updates. Configuration is loaded from environment variables.
By default all logging is printed to stdout, but a log file can be configured wih the ```-log``` commandline parameter.

### Configuration variables
* NAME_DEV_MODE - TRUE to run against the name.com dev server
* NAME_HOSTNAMES - Comma seperated list of hostnames to update
* NAME_DOMAIN - Base domain name
* NAME_INTERVAL - Interval (seconds) between updates
* NAME_TOKEN - Name.com API token
* NAME_USER - Name.com API username

## Error Handling

Currently, there is limited testing, primarily on name-api dependant utilities.
While error handling _should_ be done gracefully, not every edge case has been tested.

Ideally, when running in daemon mode, name-dyndns tries to treat any errors
arising from network as transient failures, and tries again next iteration. The idea behind this is that a single network failure shouldn't
kill the daemon, which could then potentially result in having the DNS records out
of sync, which would defeat the whole point of name-dyndns.
