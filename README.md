# go-quandl

Golang package for hitting the quandl API

## Usage

In order to use deault inputs, including a lookup of all 500 S&P500 symbols given in your go-quandl/data/inputs/SP500.json, the following command saves your input symbol(s) WIKI data by day to the data/output folder on your local machine use:

```
go run main.go -api_key=$QUANDLAPIKEY -path=$GOPATH/src/github.com/twold/go-quandl/data
```

## Output

The default option saves data using the following folder/file naming convention:

../go-quandl/data/output/<SYMBOL>/<YYYY-MM-DD>.json

Each file has a json object with the following fields

| Name 			| Type		|
|:--------------|:----------|
| Date       	| string	|
| DayOfWeek  	| string	|
| Open       	| float64	|
| High       	| float64	|
| Low        	| float64	|
| Close      	| float64	|
| Volume     	| float64	|
| ExDividend 	| float64	|
| SplitRatio 	| float64	|
| AdjOpen    	| float64	|
| AdjHigh    	| float64	|
| AdjLow     	| float64	|
| AdjClose   	| float64	|
| AdjVolume  	| float64	|


.../AAPL/1980-12-12.json

```
{
	"Date": "1980-12-12",
	"DayOfWeek": "Friday",
	"Open": 28.75,
	"High": 28.87,
	"Low": 28.75,
	"Close": 28.75,
	"Volume": 2093900,
	"Ex-Dividend": 0,
	"Split Ratio": 1,
	"Adj. Open": 0.42270591588018,
	"Adj. High": 0.42447025361603,
	"Adj. Low": 0.42270591588018,
	"Adj. Close": 0.42270591588018,
	"Adj. Volume": 117258400
}
```