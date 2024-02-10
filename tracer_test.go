package rxlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
)

func TestNewLogger(t *testing.T) {

	newTracer := NewTracer(&Opts{})
	logger, err := newTracer.NewLogger("BAC-NET", "BAC-1234", ColorBlue, 100)
	if err != nil {
		fmt.Println(err)
		return
	}

	newTrace := logger.NewTrace("NET", ColorYellow, 5)

	newTrace.Errorf("err from a trace: %d", 1234)
	newTrace.Logf("hello from a message")
	newTrace.Debugf("hello from a message")
	newTrace.Errorf("hello from a message")

	pprint.PrintJSON(logger.Traces)

}
