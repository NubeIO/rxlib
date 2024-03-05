package payload

import (
	"fmt"
	"testing"
)

func TestNewPayload(t *testing.T) {
	p, _ := NewPayload(&DataPayload{
		PortID:   "a",
		DataType: "float64",
		IsNil:    false,
		Data:     22.2,
	})
	float, err := p.ToFloat()
	if err != nil {
		return
	}
	fmt.Println(float + 10)

}

type person struct {
	Name string
}

func TestNewPayloadJSON(t *testing.T) {
	p, _ := NewPayload(&DataPayload{
		PortID:   "a",
		DataType: "json",
		IsNil:    false,
		Data:     &person{Name: "bob"},
	})

	var person *person
	p.Unmarshal(&person)
	fmt.Println(person)

}
