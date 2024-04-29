package jsonutils

import (
	"fmt"
	"testing"
)

func TestJSONUtils_Parse(t *testing.T) {
	jc := &JSONUtils{}

	j := jc.Parse(jString).Get("dayToTimeRanges")

	fmt.Println(j)
}

var jString = `
{
		"name": "test",
		"useUTC": false,
		"setTimezone": "Australia/Sydney",
		"dayToTimeRanges": {
			"0": [
				{
					"Start": "2024-04-28T07:34:00+10:00",
					"Stop": "2024-04-28T11:00:00+10:00"
				}
			],
			"3": [
				{
					"Start": "2024-04-28T19:00:00+10:00",
					"Stop": "2024-04-28T19:35:00+10:00"
				}
			]
		},
		"exceptions": {
			"2024-01-30T09:00:00Z": {
				"Start": "2024-01-30T09:00:00Z",
				"Stop": "2024-02-12T05:00:00Z"
			}
		}
	}
`
