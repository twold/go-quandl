package endpoints

import (
	"errors"
	"fmt"
)

var (
	errInvalidService = errors.New("Invalid service request.")
)

// Base Url
const (
	url    = "www.quandl.com"
	suffix = "/api"
)

// Defaults
const (
	defaultProtocol = "https"
	defaultVersion  = "v3"
)

// Service identifiers
const (
	Datasets = "datasets"
)

// Database codes
const (
	CBOE = "CBOE"
	WIKI = "WIKI"
)

// DataTypes
const (
	DATA     = "data"
	METADATA = "metadata"
)

// Return formats
const (
	JSON = "json"
	XML  = "xml"
	CSV  = "csv"
)

// Database codes
var (
	DatabaseCodes = []string{
		CBOE,
		WIKI,
	}
)

// Tickers
var (
	Formats = []string{
		JSON,
		XML,
		CSV,
	}
)

// Tickers
var (
	Types = []string{
		DATA,
		METADATA,
	}
)

var (
	Services = map[string]service{
		"datasets": service{
			Name:          Datasets,
			opts:          DatabaseCodes,
			dataTypes:     Types,
			returnFormats: Formats,
		},
	}
)

type service struct {
	Name          string
	opts          []string
	params        []string
	dataTypes     []string
	returnFormats []string
}

type Endpoint struct {
	BaseUrl         string
	DefaultProtocol string
	defaultVersion  string
	service
	suffix string
	url    string
	URL    string
}

func endpoint(name string) Endpoint {
	return Endpoint{
		BaseUrl:         url,
		DefaultProtocol: defaultProtocol,
		defaultVersion:  defaultVersion,
		service:         Services[name],
		suffix:          suffix,
		url:             url,
	}
}

func New(service, opt, param, dataType, format string) (Endpoint, error) {
	e := endpoint(service)
	if e.opts == nil {
		return e, errInvalidService
	}

	URL := fmt.Sprintf("%s://%s%s/%s/%s", e.DefaultProtocol, e.url, e.suffix, defaultVersion, service)
	for _, s := range e.opts {
		if s == opt {
			URL = fmt.Sprintf("%s/%s", URL, opt)
		}
	}

	URL = fmt.Sprintf("%s/%s", URL, param)

	for _, d := range e.dataTypes {
		if d == dataType {
			URL = fmt.Sprintf("%s/%s", URL, dataType)
		}
	}

	for _, f := range e.returnFormats {
		if f == format {
			URL = fmt.Sprintf("%s.%s", URL, format)
		}
	}
	return Endpoint{URL: URL}, nil
}
