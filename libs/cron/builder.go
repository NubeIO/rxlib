package cron

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	maxMinutes = 59
	minMinutes = 0
	maxHours   = 23
	minHours   = 0
	maxDays    = 31
	minDays    = 1
	minMonth   = time.January
	minWeekday = time.Monday
)

type monthType = time.Month

type weekdayType = time.Weekday

type cronNumber interface {
	weekdayType | monthType | int
}

type cronTime interface {
	parse() string
}

type CronSchedule struct {
	minute     cronTime
	hour       cronTime
	day        cronTime
	month      cronTime
	week       cronTime
	stringWeek string
	err        error
}

func (p *CronSchedule) New() *CronSchedule {
	return &CronSchedule{
		minute: defaultCronTime{},
		hour:   defaultCronTime{},
		day:    defaultCronTime{},
		month:  defaultCronTime{},
		week:   defaultCronTime{},
	}
}

func (p *CronSchedule) Parse() (string, error) {
	if p.err != nil {
		return "", p.err
	}
	return p.String(), nil
}

func (p *CronSchedule) setErr(err error) {
	if p.err == nil {
		p.err = err
	}
}

func (p *CronSchedule) setMinute(value cronTime) {
	p.minute = value
}

func (p *CronSchedule) setHour(value cronTime) {
	p.hour = value
}

func (p *CronSchedule) setDay(value cronTime) {
	p.day = value
}

func (p *CronSchedule) setMonth(value cronTime) {
	p.month = value
}

func (p *CronSchedule) setWeek(value cronTime) {
	p.week = value
}

func (p *CronSchedule) Minutes(minutes ...int) *CronSchedule {
	err := getThresholdError("minutes", maxMinutes, minMinutes, minutes)
	p.setErr(err)
	values := newValues(minutes)
	p.setMinute(&values)
	return p
}

func (p *CronSchedule) MinutesRange(from, to int) *CronSchedule {
	err1 := getThresholdError("minutes", maxMinutes, minMinutes, []int{from, to})
	err2 := getRangeError("minutes", from, to)
	if err1 != nil {
		p.setErr(err1)
	}
	if err2 != nil {
		p.setErr(err2)
	}
	p.setMinute(newRange(from, to))
	return p
}

func (p *CronSchedule) MinutesInterval(parameter int) *CronSchedule {
	err := getThresholdError("minutes", maxMinutes, minMinutes+1, []int{parameter})
	p.setErr(err)
	p.setMinute(newInterval(parameter))
	return p
}

func (p *CronSchedule) MinutesRangedInterval(from, to, parameter int) *CronSchedule {
	err1 := getThresholdError("minutes", from, to, []int{parameter})
	err2 := getThresholdError("minutes", math.MaxInt, minMinutes+1, []int{parameter})
	err3 := getRangeError("minutes", from, to)
	if err1 != nil {
		p.setErr(err1)
	}
	if err2 != nil {
		p.setErr(err2)
	}
	if err3 != nil {
		p.setErr(err3)
	}
	p.setMinute(newRangedInterval(from, to, parameter))
	return p
}

func (p *CronSchedule) Hours(hours ...int) *CronSchedule {
	err := getThresholdError("hours", maxHours, minHours, hours)
	p.setErr(err)
	values := newValues(hours)
	p.setHour(&values)
	return p
}

func (p *CronSchedule) HoursRange(from, to int) *CronSchedule {
	err1 := getThresholdError("hours", maxHours, minHours, []int{from, to})
	err2 := getRangeError("hours", from, to)
	if err1 != nil {
		p.setErr(err1)
	}
	if err2 != nil {
		p.setErr(err2)
	}
	p.setHour(newRange(from, to))
	return p
}

func (p *CronSchedule) HoursInterval(parameter int) *CronSchedule {
	err := getThresholdError("hours", maxHours, minHours+1, []int{parameter})
	p.setErr(err)
	p.setHour(newInterval(parameter))
	return p
}

func (p *CronSchedule) HoursRangedInterval(from, to, parameter int) *CronSchedule {
	err1 := getThresholdError("hours", maxHours, minHours, []int{from, to})
	err2 := getThresholdError("hours", maxHours, minHours+1, []int{parameter})
	err3 := getRangeError("hours", from, to)
	if err1 != nil {
		p.setErr(err1)
	}
	if err2 != nil {
		p.setErr(err2)
	}
	if err3 != nil {
		p.setErr(err3)
	}
	p.setHour(newRangedInterval(from, to, parameter))
	return p
}

func (p *CronSchedule) Day(dayName string) *CronSchedule {
	dayName = strings.ToUpper(dayName)
	dayAbbreviations := map[string]int{
		"SUN": 0,
		"MON": 1,
		"TUE": 2,
		"WED": 3,
		"THU": 4,
		"FRI": 5,
		"SAT": 6,
	}

	cronDigit, found := dayAbbreviations[dayName]
	if !found {
		p.setErr(fmt.Errorf("invalid day name: %s", dayName))
		return p
	}

	// Set both the integer-based day of the week (0-6) and the string-based day abbreviation (SUN, MON, etc.)
	p.week = newValues([]int{cronDigit})
	p.stringWeek = dayName
	return p
}

func (p *CronSchedule) String() string {
	week := "*"
	if p.stringWeek != "" {
		week = p.stringWeek
	}
	return fmt.Sprintf("%s %s %s %s %s", p.minute.parse(), p.hour.parse(), p.day.parse(), p.month.parse(), week)
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0 // Default to 0 if conversion fails
	}
	return i
}

func (p *CronSchedule) Days(days ...int) *CronSchedule {
	err := getThresholdError("days", maxDays, minDays, days)
	p.setErr(err)
	values := newValues(days)
	p.setDay(&values)
	return p
}

func (p *CronSchedule) DaysRange(from, to int) *CronSchedule {
	err1 := getThresholdError("days", maxDays, minDays, []int{from, to})
	err2 := getRangeError("days", from, to)
	if err1 != nil {
		p.setErr(err1)
	}
	if err2 != nil {
		p.setErr(err2)
	}
	p.setDay(newRange(from, to))
	return p
}

func (p *CronSchedule) DaysInterval(parameter int) *CronSchedule {
	err := getThresholdError("days", maxDays, minDays+1, []int{parameter})
	p.setErr(err)
	p.setDay(newInterval(parameter))
	return p
}

func (p *CronSchedule) DaysRangedInterval(from, to, parameter int) *CronSchedule {
	err1 := getThresholdError("days", maxDays, minDays, []int{from, to})
	err2 := getThresholdError("days", maxDays, minDays+1, []int{parameter})
	err3 := getRangeError("days", from, to)
	if err1 != nil {
		p.setErr(err1)
	}
	if err2 != nil {
		p.setErr(err2)
	}
	if err3 != nil {
		p.setErr(err3)
	}
	p.setDay(newRangedInterval(from, to, parameter))
	return p
}

func (p *CronSchedule) Months(month ...monthType) *CronSchedule {
	values := newValues(month)
	p.setMonth(&values)
	return p
}

func (p *CronSchedule) MonthsRange(from, to monthType) *CronSchedule {
	err := getRangeError("months", int(from), int(to))
	p.setErr(err)
	p.setMonth(newRange(from, to))
	return p
}

func (p *CronSchedule) MonthsInterval(parameter int) *CronSchedule {
	err := getThresholdError("months", math.MaxInt, int(minMonth), []int{parameter})
	p.setErr(err)
	p.setMonth(newInterval(parameter))
	return p
}

func (p *CronSchedule) MonthsRangedInterval(from, to monthType, parameter int) *CronSchedule {
	err1 := getThresholdError("months", math.MaxInt, int(minMonth), []int{parameter})
	err2 := getRangeError("months", int(from), int(to))
	if err1 != nil {
		p.setErr(err1)
	}
	if err2 != nil {
		p.setErr(err2)
	}
	p.setMonth(newRangedInterval(int(from), int(to), parameter))
	return p
}

func (p *CronSchedule) Weeks(weekday ...weekdayType) *CronSchedule {
	values := newValues(weekday)
	p.setWeek(&values)
	return p
}

func (p *CronSchedule) WeeksRange(from, to weekdayType) *CronSchedule {
	err := getRangeError("weeks", int(from), int(to))
	p.setErr(err)
	p.setWeek(newRange(from, to))
	return p
}

func (p *CronSchedule) WeeksInterval(parameter int) *CronSchedule {
	err := getThresholdError("weeks", math.MaxInt, int(minWeekday), []int{parameter})
	p.setErr(err)
	p.setWeek(newInterval(parameter))
	return p
}

func (p *CronSchedule) WeeksRangedInterval(from, to weekdayType, parameter int) *CronSchedule {
	err1 := getThresholdError("weeks", math.MaxInt, int(minWeekday), []int{parameter})
	err2 := getRangeError("weeks", int(from), int(to))
	if err1 != nil {
		p.setErr(err1)
	}
	if err2 != nil {
		p.setErr(err2)
	}
	p.setWeek(newRangedInterval(int(from), int(to), parameter))
	return p
}

func NewBuilder() *CronSchedule {
	return &CronSchedule{
		minute: defaultCronTime{},
		hour:   defaultCronTime{},
		day:    defaultCronTime{},
		month:  defaultCronTime{},
		week:   defaultCronTime{},
	}
}

type CronBuilder struct {
	Day     string
	Minutes *int
	Hours   *int
	Days    *int
	Months  *int
}

// Builder make an expression by passing in the &CronBuilder{Day:"MON", Minutes:0, Hours:6} -> 0 6 * * MON
func (p *CronSchedule) Builder(c *CronBuilder) string {
	if c.Day != "" {
		p.Day(c.Day)
	}
	if c.Minutes != nil {
		p.Minutes(*c.Minutes)
	}
	if c.Hours != nil {
		p.Hours(*c.Hours)
	}
	if c.Months != nil {
		p.Months(time.Month(*c.Months))
	}
	if c.Days != nil {
		p.Days(*c.Days)
	}

	return p.String()
}

type cronValues[T cronNumber] []T

func (v cronValues[T]) isValid() bool {
	return len(v) > 0
}

func (v cronValues[T]) parse() string {
	if !v.isValid() {
		return "*"
	}
	stringValues := make([]string, len(v))
	for i, value := range v {
		stringValues[i] = fmt.Sprintf("%d", value)
	}
	return strings.Join(stringValues, ",")
}

type cronRange[T cronNumber] struct {
	from T
	to   T
}

func (r cronRange[T]) isValid() bool {
	return r.from > 0 || r.to > 0 && r.to > r.from
}

func (r cronRange[T]) parse() string {
	if !r.isValid() {
		return "*"
	}
	return fmt.Sprintf("%d-%d", r.from, r.to)
}

type cronInterval[T cronNumber] struct {
	parameter T
}

func (i cronInterval[T]) isValid() bool {
	return i.parameter > 1
}

func (i cronInterval[T]) parse() string {
	if !i.isValid() {
		return "*"
	}
	parameter := i.parameter
	return fmt.Sprintf("*/%d", parameter)
}

type cronRangedInterval[T cronNumber] struct {
	molecule  cronRange[T]
	parameter T
}

func (i cronRangedInterval[T]) isValid() bool {
	return i.molecule.isValid() && i.parameter > 1
}

func (i cronRangedInterval[T]) parse() string {
	if !i.isValid() {
		return "*"
	}
	molecule := i.molecule
	parameter := i.parameter
	moleculeStr := molecule.parse()
	return fmt.Sprintf("%s/%d", moleculeStr, parameter)
}

type defaultCronTime struct {
}

func (d defaultCronTime) parse() string {
	return "*"
}

func newValues[T cronNumber](v []T) cronValues[T] {
	return append([]T{}, v...)
}

func newRange[T cronNumber](from, to T) *cronRange[T] {
	ranges := &cronRange[T]{
		from, to,
	}
	return ranges
}

func newInterval[T cronNumber](parameter T) *cronInterval[T] {
	ranges := &cronInterval[T]{
		parameter,
	}
	return ranges
}

func newRangedInterval[T cronNumber](from, to, parameter T) *cronRangedInterval[T] {
	ranges := newRange(from, to)
	interval := &cronRangedInterval[T]{
		molecule:  *ranges,
		parameter: parameter,
	}
	return interval
}
