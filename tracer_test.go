package rxlib

import (
	"fmt"
	"testing"
)

func TestNewLogger(t *testing.T) {

	newTracer := NewTracer(&Opts{})
	logger, err := newTracer.NewLogger("BAC-NET", "BAC-1234", ColorBlue, 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	newTrace := logger.NewTrace("NET", ColorYellow, 5)

	newTrace.Errorf("err from a trace: %d", 1234)
	newTrace.Logf("hello from a message")

}
