package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestCommandBuilder(t *testing.T) {

	cmdString := `writeInput --query:(objects:name == math-add-2) --id:"in 1" --field:"value" --value:22.5` // []string
	cmdString = `setObject --name:"math-add-2" --field:"name" --value:"new name"`                           // string
	cmdString = `getObject --name:"math-add-2"`                                                             // object
	cmdString = `getObjects --query:(objects:name == math-add-2)`                                           // []Object
	cmdString = `get inputs --query="(objects:name == math-add-2)"  --field=123  --f=23 --a=2 --ddd=222 --aoo=555`
	//cmdString = `getObjects --return:json`
	cmd := NewCommand()
	cmd, err := cmd.Parse(cmdString)
	if err != nil {
		return
	}
	pprint.PrintJSON(cmd)

	fmt.Println(cmd.GetArgsByIndex(0), cmd.GetArgsByIndex(1))

}
