package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/twold/go-quandl/api"
)

var (
	api_key   string
	datatype  string
	dbcode    string
	format    string
	inputFile string
	sector    string
	path      string
	ticker    string

	err     error
	tickers []string
)

func init() {
	// Add api key
	flag.StringVar(&api_key, "api_key", "$QUANDLAPIKEY", "-api_key=xyzABCD1234567890 add api key for data pull")
	// 'data' is all we need right now
	flag.StringVar(&datatype, "datatype", "data", "-datatype=metadata your two options are 'data' and 'metadata'.  Default is 'data'")
	// WIKI is fully baked, need help with CB
	flag.StringVar(&dbcode, "dbcode", "WIKI", "-dbcode=WIKI add dbcode you would like to query here")
	// json option is fully baked
	flag.StringVar(&format, "format", "json", "-format=json add the response format you would like here.  Options are 'json', xml' and 'csv'")
	// Modified SP500 input file from
	// https://pkgstore.datahub.io/core/s-and-p-500-companies/constituents_json/data/64dd3e9582b936b0352fdd826ecd3c95/constituents_json.json
	flag.StringVar(&inputFile, "inputFile", "constituents_json.json", "-inputFile=")
	// Select sector
	flag.StringVar(&sector, "sector", "all", "-sector=all returns all symbols in list, or you can specify a sector to retrieve a subset of symbols 'Industrials', 'Health Care', 'Information Technology'")
	// This is path to data files including input and output folders
	flag.StringVar(&path, "path", "", "-path=")
	// Optional input param.  This overwrites input file with multiple symbols if given.
	flag.StringVar(&ticker, "ticker", "", "-ticker=FB input ticker symbol to retrieve data set")
}

// sample input where $QUANDLAPIKEY is your api key and $GOPATH/src/github.com/twold/go-quandl/data
// is where you have input file and is desired output location

// go run main.go -api_key=$QUANDLAPIKEY -path=$GOPATH/src/github.com/twold/go-quandl/data

func main() {
	flag.Parse()

	// create service using input format and data type
	svc := api.New(&datatype, &dbcode, &format, &api_key)

	// if individual ticker input is not given, read input file
	if ticker == "" {
		tickers, err = api.ReadInputList(path, inputFile, sector)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		tickers = append(tickers, ticker)
	}

	// range over inputs and retrieve data from API
	for _, ticker := range tickers {
		// pull data from API
		resp, err := svc.Get(path, ticker)
		if err != nil {
			// ignore invlaid ticker properly formatted eerror message from quandl
			if strings.Contains(err.Error(), "Quandl code. Please check your Quandl codes and try again") {
				log.Printf("Remove ticker from list %+v.\n Error ignored: %+v\n", ticker, err)
				continue
			}
			// ignore URL parsing errors due to invalid symbol
			if strings.Contains(err.Error(), "We could not recognize the URL you requested") {
				log.Printf("Remove ticker from list %+v.\n Error ignored: %+v\n", ticker, err)
				continue
			}
			// Print error to std out and exit
			log.Fatalln(err)
		}

		// response will be null if you chose "metadata" as your datatype
		if resp != nil {

			// ToDo: need switch on xml/json
			// format json response
			obj, err := json.MarshalIndent(resp.Data, "", "	")
			if err != nil {
				log.Fatal(err)
			}
			// output to stdout
			fmt.Sprintf("%s\n", string(obj))
		}
	}
}
