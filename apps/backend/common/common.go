package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseDates(str string) string {

	re := regexp.MustCompile("::(year|month|day):(\\+|-)?(\\d+)?:")

	matches := re.FindAllStringSubmatch(str, -1)
	now := time.Now()
	for _, match := range matches {

		yearVal := 0
		monthVal := 0
		dayVal := 0
		if len(match) == 4 {
			val := 0
			if v, err := strconv.Atoi(match[3]); err == nil {
				val = v
			}
			if match[2] == "-" {
				val *= -1
			}
			if match[1] == "year" {

				yearVal = val
				str = strings.ReplaceAll(str, match[0], "#year#")
			} else if match[1] == "month" {
				monthVal = val
				str = strings.ReplaceAll(str, match[0], "#month#")
			} else if match[1] == "day" {
				dayVal = val
				str = strings.ReplaceAll(str, match[0], "#day#")
			}
		}
		now = now.AddDate(yearVal, monthVal, dayVal)
		str = strings.ReplaceAll(str, "#year#", fmt.Sprintf("%d", now.Year()))
		str = strings.ReplaceAll(str, "#month#", fmt.Sprintf("%d", now.Month()))
		str = strings.ReplaceAll(str, "#day#", fmt.Sprintf("%d", now.Day()))

	}

	re2 := regexp.MustCompile("::monthdays::")

	matches2 := re2.FindAllStringSubmatch(str, -1)
	dinm := daysInMonth(now)
	for _, match := range matches2 {
		str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%d", dinm))
	}

	return str
}

func daysInMonth(t time.Time) int {
	t = time.Date(t.Year(), t.Month(), 32, 0, 0, 0, 0, time.UTC)
	daysInMonth := 32 - t.Day()
	days := make([]int, daysInMonth)
	for i := range days {
		days[i] = i + 1
	}

	d := days[len(days)-1]
	d += 1
	return d
}
