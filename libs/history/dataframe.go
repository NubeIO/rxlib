package history

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"sync"
	"time"
)

type DataFrameOperations interface {
	Avg() float64
	Sum() float64
	Min(columnName string) float64
	Max(columnName string) float64
	FilterData(comparator series.Comparator, value interface{}, columnName string) DataFrameOperations
	FilterDataBetween(startDate, endDate time.Time, columnName string) DataFrameOperations
	FilterDateRange(startDate, endDate string, columnName string) DataFrameOperations
	FilterByTimeDays(startTime, endTime string, days ...int) DataFrameOperations
	FilterByDateDays(startDate, endDate string, days ...int) DataFrameOperations
	ToDF() dataframe.DataFrame
	SetDF(df dataframe.DataFrame)
	Count() int
}

type data struct {
	df   dataframe.DataFrame
	lock sync.Mutex
}

func (d *data) Avg() float64 {
	d.lock.Lock()
	defer d.lock.Unlock()
	meanVal := d.df.Col("Value").Mean()
	return meanVal
}

func (d *data) Sum() float64 {
	d.lock.Lock()
	defer d.lock.Unlock()
	meanVal := d.df.Col("Value").Sum()
	return meanVal
}

func (d *data) Count() int {
	d.lock.Lock()
	defer d.lock.Unlock()
	meanVal := d.df.Col("Value")
	return meanVal.Len()
}

func (d *data) Min(columnName string) float64 {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.df.Col(columnName).Min()
}

func (d *data) Max(columnName string) float64 {
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

func (d *data) ToDF() dataframe.DataFrame {
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

func (d *data) FilterByTimeDays(startTime, endTime string, days ...int) DataFrameOperations {
	d.lock.Lock()
	defer d.lock.Unlock()

	startTimeParsed, errStart := time.Parse("15:04", startTime)
	if errStart != nil {
		fmt.Println("Error parsing startTime:", errStart)
		return d
	}
	endTimeParsed, errEnd := time.Parse("15:04", endTime)
	if errEnd != nil {
		fmt.Println("Error parsing endTime:", errEnd)
		return d
	}

	dayMap := make(map[int]bool)
	for _, day := range days {
		dayMap[day] = true
	}

	filterFunc := func(el series.Element) bool {
		ts, ok := el.Val().(string)
		if !ok {
			fmt.Println("Invalid type for Timestamp")
			return false
		}
		timestamp, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			fmt.Println("Error parsing Timestamp:", ts)
			return false
		}

		dayOfWeek := int(timestamp.Weekday())
		timeOfDay := timestamp.Format("15:04")
		timeOfDayParsed, _ := time.Parse("15:04", timeOfDay)

		match := dayMap[dayOfWeek] && !timeOfDayParsed.Before(startTimeParsed) && !timeOfDayParsed.After(endTimeParsed)
		fmt.Printf("Testing Timestamp: %s, Day of the Week: %s, Match: %t\n", timestamp, timestamp.Weekday().String(), match)
		return match
	}

	filteredDf := d.df.Filter(dataframe.F{
		Colname:    "Timestamp",
		Comparator: series.CompFunc,
		Comparando: filterFunc,
	})

	return &data{df: filteredDf}
}

func (d *data) FilterByDateDays(startDate, endDate string, days ...int) DataFrameOperations {
	d.lock.Lock()
	defer d.lock.Unlock()

	// Parse the start and end date and times
	startDateTime, errStart := time.Parse("2006-01-02:15:04", startDate)
	if errStart != nil {
		fmt.Println("Error parsing start date and time:", errStart)
		return d
	}
	endDateTime, errEnd := time.Parse("2006-01-02:15:04", endDate)
	if errEnd != nil {
		fmt.Println("Error parsing end date and time:", errEnd)
		return d
	}

	// Create map for quick day lookup
	dayMap := make(map[int]bool)
	for _, day := range days {
		dayMap[day] = true
	}

	// Define the filter function
	filterFunc := func(el series.Element) bool {
		ts, ok := el.Val().(string)
		if !ok {
			return false
		}
		timestamp, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return false
		}

		dayOfWeek := int(timestamp.Weekday())
		return dayMap[dayOfWeek] && !timestamp.Before(startDateTime) && !timestamp.After(endDateTime)
	}

	// Apply the filter
	filteredDf := d.df.Filter(dataframe.F{
		Colname:    "Timestamp",
		Comparator: series.CompFunc,
		Comparando: filterFunc,
	})

	return &data{df: filteredDf}
}
