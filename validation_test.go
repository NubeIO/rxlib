package rxlib

import (
	"fmt"
	"testing"
)

func TestNewObjectBuilder(t *testing.T) {
	key := "key"
	builder := NewObjectBuilder().AddValidation(key)
	validation := builder.GetValidation(key)
	if validation != nil {
		validation.SetError(fmt.Errorf("Sample error message"))
	}
	fmt.Printf("%+v\n", validation)
}
