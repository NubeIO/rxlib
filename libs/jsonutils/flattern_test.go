package jsonutils

import (
	"github.com/NubeIO/rxlib/helpers/pprint"
	"github.com/tidwall/gjson"
	"testing"
)

func TestFlattenJSON(t *testing.T) {
	jsonStr := `{
   	"isbn": "123-456-222",
   	"author": {
   	  "lastname": "Doe",
   	  "firstname": "Jane"
   	},
   	"editor": {
   	  "lastname": "Smith",
   	  "firstname": "Jane"
   	},
   	"title": "The Ultimate Database Study Guide",
   	"category": [
   	  "Non-Fiction",
   	  "Technology"
   	]
     }`

	f, _ := FlattenJSON(jsonStr, DotSeparator)

	p := gjson.Parse(f)

	pprint.PrintJSON(p.Value())
}
