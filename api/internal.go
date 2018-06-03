package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/xitongsys/parquet-go/ParquetFile"
	"github.com/xitongsys/parquet-go/ParquetWriter"
)

func formatDataSet(columnNames []*string, data []interface{}) ([]byte, error) {
	log.Printf("Transforming data set.\n")
	str := "["
	for n, objs := range data {
		str = fmt.Sprintf("%s{", str)
		typ := reflect.TypeOf(objs)
		if typ.Kind() == reflect.Slice {
			for i, obj := range objs.([]interface{}) {
				name := columnNames[i]
				if *name == "Date" {
					date, err := time.Parse("2006-01-02", obj.(string))
					if err != nil {
						return nil, err
					}
					str = fmt.Sprintf("%s\"DayOfWeek\": \"%v\",", str, date.Weekday())
				}
				obj := obj
				typ = reflect.TypeOf(obj)
				switch typ.Kind() {

				// format float64 field with column names
				case reflect.Float64:
					str = fmt.Sprintf("%s\"%s\": %g", str, *name, obj.(float64))

				// format string field with column names
				case reflect.String:
					str = fmt.Sprintf("%s\"%s\": \"%v\"", str, *name, obj.(string))

				// format unknown types as string
				default:
					str = fmt.Sprintf("%s\"%s\": \"%v\"", str, *name, obj)
				}
				// add comma between fields
				if i < (len(objs.([]interface{})) - 1) {
					str = fmt.Sprintf("%s,", str)
				}
			}
		}

		str = fmt.Sprintf("%s}", str)
		// add comma between fields
		if n < (len(data) - 1) {
			str = fmt.Sprintf("%s,", str)
		}
	}
	str = fmt.Sprintf("%s]", str)
	return []byte(str), nil
}

func read(body io.ReadCloser) ([]byte, error) {
	log.Printf("Reading API response.\n")
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalInputFile(data []byte) ([]List, error) {
	var dat []List
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (c *CBOE) unmarshal(data []byte) (*CBOE, error) {
	var dat CBOE
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, nil
}

func (c *CBOE) unmarshalData(data []byte) ([]CBOE, error) {
	var dat []CBOE
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (c *Wiki) unmarshal(data []byte) (*Wiki, error) {
	var dat Wiki
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, nil
}

func (c *Wiki) unmarshalData(data []byte) ([]Wiki, error) {
	var dat []Wiki
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

type Err struct {
	QuandlError `json:"quandl_error,omitempty" type:"struct"`
}

type QuandlError struct {
	Code    *string `json:"code,omitempty" type:"string"`
	Message *string `json:"message" type:"string"`
}

func unmarshalError(data []byte) (*Err, error) {
	var dat Err
	err := json.Unmarshal(data, &dat)
	if err != nil {
		return nil, err
	}
	return &dat, nil
}

func readFile(path, fileName string) ([]byte, error) {

	b, err := ioutil.ReadFile(filepath.Join(path, "input", fileName))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func writeToFiles(path, symbol string, objs []Wiki) error {
	err := os.Mkdir(filepath.Join(path, "output", symbol), 0777)
	if err != nil {
		if check := strings.Contains(err.Error(), "file exists"); check == false {
			return err
		}
	}

	for _, obj := range objs {

		name := filepath.Join(path, "output", symbol, fmt.Sprintf("%v.txt", *obj.Date))
		if _, err := os.Stat(name); os.IsNotExist(err) == false {
			continue
		}

		fw, err := ParquetFile.NewLocalFileWriter(name)
		if err != nil {
			log.Println("Can't create file", err)
			return err
		}

		//write
		pw, err := ParquetWriter.NewParquetWriter(fw, new(Wiki), 4)
		if err != nil {
			log.Println("Can't create parquet writer", err)
			return err
		}

		if err = pw.Write(obj); err != nil {
			log.Println("Write error", err)
			return err
		}

		if err = pw.WriteStop(); err != nil {
			log.Println("WriteStop error", err)
			return err
		}
		log.Printf("Write Finished for file %+v\n", name)
		fw.Close()

	}
	return nil
}

func writeLocalFiles(path, symbol string, objs []Wiki) error {
	log.Printf("Writing local files.\n")
	err := os.Mkdir(filepath.Join(path, "output", symbol), 0777)
	if err != nil {
		if check := strings.Contains(err.Error(), "file exists"); check == false {
			return err
		}
	}

	for _, obj := range objs {
		b, err := json.MarshalIndent(obj, "", "	")
		if err != nil {
			return err
		}

		name := filepath.Join(path, "output", symbol, fmt.Sprintf("%v.json", *obj.Date))
		if _, err := os.Stat(name); os.IsNotExist(err) == false {
			continue
		}

		f, err := os.Create(name)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}
