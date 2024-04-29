package topics

import (
	"fmt"
	"testing"
)

func TestNewPublishTopicBuilder(t *testing.T) {
	pubBuilder := NewPublishTopicBuilder().
		WithVersion("v1").
		WithClientType(Cloud).
		WithSenderUUID("R-1").
		WithTargetUUID("target-123").
		WithRequestUUID("request-789").
		WithDataType(Proto).
		WithType(Object)
	fmt.Println(pubBuilder.Build())

	// r/req/v1/cloud/R-1/plain/command/RX-2/req-uuid

}
