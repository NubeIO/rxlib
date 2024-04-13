package history

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"sync"
	"time"
)

type DataFrameOperations interface {
	CalculateAverage() float64
	CalculateMin(columnName string) float64
	CalculateMax(columnName string) float64
	FilterData(comparator series.Comparator, value interface{}, columnName string) DataFrameOperations
	FilterDataBetween(startDate, endDate time.Time, columnName string) DataFrameOperations
	FilterDateRange(startDate, endDate string, columnName string) DataFrameOperations
	GetDF() dataframe.DataFrame
	SetDF(df dataframe.DataFrame)
}

type data struct {
	df   dataframe.DataFrame
	lock sync.Mutex
}

func (d *data) CalculateAverage() float64 {
	d.lock.Lock()
	defer d.lock.Unlock()
	meanVal := d.df.Col("Value").Mean()
	return meanVal
}

func (d *data) CalculateMin(columnName string) float64 {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.df.Col(columnName).Min()
}

func (d *data) CalculateMax(columnName string) float64 {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.df.Col(columnName).Max()
}

func (d *data) FilterData(comparator series.Comparator, value interface{}, columnName string) DataFrameOperations {
	d.lock.Lock()
	defer d.lock.Unlock()
	filteredDf := d.df.Filter(dataframe.F{
		Colname:    columnName,
		Comparator: comparator,
		Comparando: value,
	})
	return &data{df: filteredDf}
}

func (d *data) FilterDateRange(startDate, endDate string, columnName string) DataFrameOperations {
	d.lock.Lock()
	defer d.lock.Unlock()

	// Parse the start and end date strings
	start, errStart := ParseDateTime(startDate)
	if errStart != nil {
		return d
	}
	end, errEnd := ParseDateTime(endDate)
	if errEnd != nil {
		return d
	}

	// Apply the filtering based on the parsed dates
	filteredDf := d.df.Filter(dataframe.F{
		Colname:    columnName,
		Comparator: series.GreaterEq,
		Comparando: start.Format(time.RFC3339),
	}).Filter(dataframe.F{
		Colname:    columnName,
		Comparator: series.LessEq,
		Comparando: end.Format(time.RFC3339),
	})
	return &data{df: filteredDf}
}

func (d *data) FilterDataBetween(startDate, endDate time.Time, columnName string) DataFrameOperations {
	d.lock.Lock()
	defer d.lock.Unlock()
	filteredDf := d.df.Filter(dataframe.F{
		Colname:    columnName,
		Comparator: series.GreaterEq,
		Comparando: startDate.Format(time.RFC3339),
	}).Filter(dataframe.F{
		Colname:    columnName,
		Comparator: series.Less,
		Comparando: endDate.Format(time.RFC3339),
	})
	return &data{df: filteredDf}
}

func (d *data) GetDF() dataframe.DataFrame {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.df
}

func (d *data) SetDF(df dataframe.DataFrame) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.df = df
}

func New(histories []*AllHistories) DataFrameOperations {
	uuids := make([]string, 0)
	values := make([]float64, 0)
	timestamps := make([]string, 0)

	for _, history := range histories {
		for _, record := range history.Histories {
			uuids = append(uuids, history.ObjectUUID)
			values = append(values, record.GetValue().(float64))
			timestamps = append(timestamps, record.GetTimestamp().Format(time.RFC3339))
		}
	}

	df := dataframe.New(
		series.New(uuids, series.String, "ObjectUUID"),
		series.New(values, series.Float, "Value"),
		series.New(timestamps, series.String, "Timestamp"),
	)
	return &data{df: df}
}

const CustomDateTimeFormat = "2006-01-02:15:04"

// ParseDateTime takes a datetime string in custom format and returns the parsed time.Time object
func ParseDateTime(datetimeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(CustomDateTimeFormat, datetimeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid datetime format: %v", err)
	}
	return parsedTime, nil
}
