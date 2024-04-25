package ndocs

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestNew(t *testing.T) {
	docs := New(ObjectString)
	a := cleanSearchTerm("runtime.GetInp().GetName()")
	fuzzy := docs.Fuzzy(a)

	pprint.PrintJSON(fuzzy)
	fmt.Println(a)
}
