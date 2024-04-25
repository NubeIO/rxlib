package main

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/rxlib/libs/docgen"
	"os"
)

func main() {

	filename := "/home/aidan/code/go/rxlib/runtime.go"
	interfaceName := "Runtime"

	helpers, err := docgen.ParseInterfaceMethods(filename, interfaceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		return
	}

	jsonData, err := json.MarshalIndent(helpers, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}
