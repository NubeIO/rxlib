package rxlib

import (
	"fmt"
	"testing"
)

func TestNewObjectBuilder(t *testing.T) {
	key := "key"
	builder := NewObjectBuilder().AddValidation(key)
	validation, _ := builder.GetValidation(key)
	if validation != nil {
		validation.SetHaltReason(&ValidationMessage{
			//Error:       fmt.Errorf("ERROROR"),
			Message:     "Message Message",
			Explanation: "Explanation Explanation",
		})

	}

	dump, _ := builder.GetValidation(key)
	fmt.Println(dump.Halt.Error)

	fmt.Printf("%+v\n", validation.ToString())
}

//func PrintJOSN(x interface{}) {
//	ioWriter := os.Stdout
//	w := json.NewEncoder(ioWriter)
//	w.SetIndent("", "    ")
//	w.Encode(x)
//}
