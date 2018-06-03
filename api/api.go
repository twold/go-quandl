package api

import (
	"errors"

	"github.com/twold/go-quandl/client"
)

type InputFile []List

type List struct {
	Name   string
	Sector string
	Symbol string
}

func ReadInputList(path, inputFile, sector string) ([]string, error) {
	b, err := readFile(path, inputFile)
	if err != nil {
		return nil, err
	}
	list, err := unmarshalInputFile(b)
	if err != nil {
		return nil, err
	}

	symbols := make([]string, 0)
	for _, item := range list {
		// if sector is not specified or if all sectors are requested
		if sector == "" || sector == "all" {
			symbols = append(symbols, item.Symbol)
			continue
		}
		// match if sector is specified
		if sector == item.Sector {
			symbols = append(symbols, item.Symbol)
		}
	}
	return symbols, nil
}

type Service struct {
	*client.Client

	dataType *string

	dbCode *string

	format *string

	DataSetData DataSet `json:"dataset_data" type:"struct"` // WIKI uses dataset_data

	DataSet `json:"dataset" type:"struct"` // CBOE uses dataset
}

type DataSet struct {
	ColumnIndex *string `json:"column_index" type:"string"`

	ColumnNames []*string `json:"column_names" type:"list"`

	RawData []interface{} `json:"data" type:"list"`

	Data interface{} `json:"Data" type:"list"`

	EndDate *string `json:"end_date" type:"string"`

	Frequency *string `json:"frequency" type:"string"`

	Limit *string `json:"limit" type:"string"`

	NewestAvailableDate *string `json:"newest_available_date" type:"string"`

	OldestAvailableDate *string `json:"oldest_available_date" type:"string"`

	StartDate *string `json:"start_date" type:"string"`

	Transform *string `json:"transform" type:"string"`
}

type Getter interface {
	Get(path, symbol string) (*DataSet, error)
}

type CBOE struct {
	*Service

	TradeDate           *string  `json:"Trade Date" type:"string"`
	DayOfWeek           *string  `json:"DayOfWeek" type:"string"`
	Open                *float64 `json:"Open" type:"float64"`
	High                *float64 `json:"High" type:"float64"`
	Low                 *float64 `json:"Low" type:"float64"`
	Close               *float64 `json:"Close" type:"float64"`
	Settle              *float64 `json:"Settle" type:"float64"`
	Change              *float64 `json:"Change" type:"float64"`
	TotalVolume         *float64 `json:"Total Volume" type:"float64"`
	EFP                 *float64 `json:"EFP" type:"float64"`
	PrevDayOpenInterest *float64 `json:"Prev. Day Open Interest" type:"float64"`
}

type Wiki struct {
	*Service

	Date       *string  `json:"Date" type:"string" parquet:"name=Date, inname=Date, type=UTF8, repetitiontype=OPTIONAL"`
	DayOfWeek  *string  `json:"DayOfWeek" type:"string" parquet:"name=DayOfWeek, inname=DayOfWeek, type=UTF8, repetitiontype=OPTIONAL"`
	Open       *float64 `json:"Open" type:"float64" parquet:"name=Open, innname=Open, type=DOUBLE, repetitiontype=OPTIONAL"`
	High       *float64 `json:"High" type:"float64" parquet:"name=High, inname=High, type=DOUBLE, repetitiontype=OPTIONAL"`
	Low        *float64 `json:"Low" type:"float64" parquet:"name=Low, inname=Low, type=DOUBLE, repetitiontype=OPTIONAL"`
	Close      *float64 `json:"Close" type:"float64" parquet:"name=Close, inname=Close, type=DOUBLE, repetitiontype=OPTIONAL"`
	Volume     *float64 `json:"Volume" type:"float64" parquet:"name=Volume, inname=Volume, type=DOUBLE, repetitiontype=OPTIONAL"`
	ExDividend *float64 `json:"Ex-Dividend" type:"float64" parquet:"name=ExDividend, inname=ExDividend, type=DOUBLE, repetitiontype=OPTIONAL"`
	SplitRatio *float64 `json:"Split Ratio" type:"float64" parquet:"name=SplitRatio, inname=SplitRatio, type=DOUBLE, repetitiontype=OPTIONAL"`
	AdjOpen    *float64 `json:"Adj. Open" type:"float64" parquet:"name=AdjOpen, inname=AdjOpen, type=DOUBLE, repetitiontype=OPTIONAL"`
	AdjHigh    *float64 `json:"Adj. High" type:"float64" parquet:"name=AdjHigh, iname=AdjHigh,  type=DOUBLE, repetitiontype=OPTIONAL"`
	AdjLow     *float64 `json:"Adj. Low" type:"float64" parquet:"name=AdjLow, inname=AdjLow, type=DOUBLE, repetitiontype=OPTIONAL"`
	AdjClose   *float64 `json:"Adj. Close" type:"float64" parquet:"name=AdjClose, inname=AdjClose, type=DOUBLE, repetitiontype=OPTIONAL"`
	AdjVolume  *float64 `json:"Adj. Volume" type:"float64" parquet:"name=AdjVolume, inname=AdjVolume, type=DOUBLE, repetitiontype=OPTIONAL"`
}

// dataType options are "data" and "metadata"
// format options are "csv", "json" and "xml"
func New(dataType, dbCode, format, key *string) Getter {

	// save all input values to client
	svc := &Service{
		dataType: dataType,
		dbCode:   dbCode,
		format:   format,
	}

	// set default to json
	if svc.format == nil {
		def := "json"
		svc.format = &def
	}

	// set default to WIKI
	if svc.dbCode == nil {
		def := "WIKI"
		svc.dbCode = &def
	}

	switch *svc.dbCode {
	case "WIKI":
		// create client
		return &Wiki{
			Service: &Service{
				Client: client.New("datasets").
					Auth(key).
					DBCode(*svc.dbCode).
					DataType(*svc.dataType).
					Format(*svc.format),
			},
		}
	}

	return &CBOE{
		Service: &Service{
			Client: client.New("datasets").
				Auth(key).
				DBCode(*svc.dbCode).
				Format(*svc.format),
		},
	}
}

func (c *Wiki) Get(path, symbol string) (*DataSet, error) {
	resp, err := c.Do("GET", symbol)
	if err != nil {
		return nil, err
	}

	b, err := read(resp.Body)
	if err != nil {
		return nil, err
	}

	// handle errors
	errQ, err := unmarshalError(b)
	if err != nil {
		return nil, err
	}

	// handle Quandl specific errors
	if errQ.Message != nil {
		return nil, errors.New(*errQ.Message)
	}

	// unmarshal struct to seperate data from API metadata
	c, err = c.unmarshal(b)
	if err != nil {
		return nil, err
	}

	// save updated struct as byte slice
	// ensure that timeseries data is indexed properly and field names are added
	byt, err := formatDataSet(c.DataSetData.ColumnNames, c.DataSetData.RawData)
	if err != nil {
		return nil, err
	}

	// format struct and return as typed slice
	d, err := c.unmarshalData(byt)
	if err != nil {
		return nil, err
	}

	err = writeLocalFiles(path, symbol, d)
	if err != nil {
		return nil, err
	}

	c.DataSetData.Data = d
	return &c.DataSetData, nil
}

// unfinished
func (c *CBOE) Get(path, symbol string) (*DataSet, error) {
	resp, err := c.Do("GET", symbol)
	if err != nil {
		return nil, err
	}

	b, err := read(resp.Body)
	if err != nil {
		return nil, err
	}

	// handle errors
	errQ, err := unmarshalError(b)
	if err != nil {
		return nil, err
	}

	// handle Quandl specific errors
	if errQ.Message != nil {
		return nil, errors.New(*errQ.Message)
	}

	// unmarshal struct to seperate data from API metadata
	c, err = c.unmarshal(b)
	if err != nil {
		return nil, err
	}

	byt, err := formatDataSet(c.DataSet.ColumnNames, c.DataSet.RawData)
	if err != nil {
		return nil, err
	}

	d, err := c.unmarshalData(byt)
	if err != nil {
		return nil, err
	}
	c.DataSetData.Data = d
	return &c.DataSetData, nil
}
