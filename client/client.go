package client

import (
	"fmt"

	"github.com/twold/go-quandl/endpoints"
	"github.com/twold/go-quandl/request"
)

type ClientInfo struct {
	APIVersion  string
	APIKey      *string
	dataType    string
	dbCode      string
	format      string
	serviceName string
}

type Client struct {
	ClientInfo
	*request.Request
}

// only service supported is "datasets"
func New(service string) *Client {
	return &Client{
		ClientInfo: ClientInfo{
			serviceName: service,
			APIVersion:  "v3",
		},
	}
}

func (c *Client) Auth(key *string) *Client {
	c.APIKey = key
	return c
}

// options are "data" and "metadata"
func (c *Client) DataType(dataType string) *Client {
	c.dataType = dataType
	return c
}

// options are "WIKI" and "EOD"
func (c *Client) DBCode(dbCode string) *Client {
	c.dbCode = dbCode
	return c
}

// options are 	"json" "xml" and "csv"
func (c *Client) Format(format string) *Client {
	c.format = format
	return c
}

func (c *Client) Do(method, ticker string) (*Client, error) {
	url, err := endpoints.New(c.serviceName, c.dbCode, ticker, c.dataType, c.format)
	if err != nil {
		return nil, err
	}
	// add API key to query string
	if c.APIKey != nil {
		url.URL = fmt.Sprintf("%s?api_key=%s", url.URL, *c.APIKey)
	}

	c.Request = request.New(method, url.URL, nil)
	if c.Error != nil {
		return nil, c.Error
	}
	return c, nil
}
