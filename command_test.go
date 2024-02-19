package rxlib

import (
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestCommandBuilder(t *testing.T) {

	cmdString := `writeInput --query:(objects:name == math-add-2) --id:"in 1" --field:"value" --value:22.5` // []string
	cmdString = `setObject --name:"math-add-2" --field:"name" --value:"new name"`                           // string
	cmdString = `getObject --name:"math-add-2"`                                                             // object
	cmdString = `getObjects --query:(objects:name == math-add-2)`                                           // []Object
	cmdString = `getInputs--query:(objects:name == math-add-2)`                                             // map[string][]*Ports
	//cmdString = `getObjects --return:json`

	cp := NewCommandParse()
	cmd, _ := cp.Parse(cmdString)
	pprint.PrintJSON(cmd)

}
