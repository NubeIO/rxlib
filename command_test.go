package rxlib

import (
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestCommandBuilder(t *testing.T) {

	cmdString := `writeInput -query:(objects:name == math-add-2) -id:"in 1" -write:22.5`
	cmdString = `getObjects -return:json`

	cp := NewCommandParse()
	cmd, _ := cp.Parse(cmdString)
	pprint.PrintJSON(cmd)

}
