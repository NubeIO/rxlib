package payload

import (
	"fmt"
	"testing"
)

func TestNewPayload(t *testing.T) {
	var f float64 = 10
	p, _ := NewPayload(&Body{
		PortID:   "a",
		DataType: "float64",
		IsNil:    false,
		Data:     f,
	})
	float, err := p.ToFloat()
	if err != nil {
		return
	}
	fmt.Println(float + 10)
	f = 20
	p.ApplyData(f)
	float, err = p.ToFloat()
	if err != nil {
		return
	}
	fmt.Println(float + 10)

}

type person struct {
	Name string
}

func TestNewPayloadJSON(t *testing.T) {
	p, _ := NewPayload(&Body{
		PortID:   "a",
		DataType: "json",
		IsNil:    false,
		Data:     &person{Name: "bob"},
	})

	var person *person
	p.Unmarshal(&person)
	fmt.Println(person)

}
