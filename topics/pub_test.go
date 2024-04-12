package topics

import (
	"fmt"
	"testing"
)

func TestNewPublishTopicBuilder(t *testing.T) {
	pubBuilder := NewPublishTopicBuilder().
		WithVersion("v1").
		WithClientType(Cloud).
		WithTargetUUID("target-123").
		WithSenderUUID("sender-456").
		WithRequestUUID("request-789").
		WithDataType(Proto).
		WithType(Object)
	fmt.Println(pubBuilder.Build())

	subBuilder := NewSubscribeTopicBuilder().
		WithVersion("v1").
		WithClientType(Cloud).
		WithTargetUUID("target-123").
		WithSenderUUID("sender-456").
		WithRequestUUID("request-789").
		WithDataType(Proto).
		WithType(Object)
	fmt.Println(subBuilder.Build())
}
