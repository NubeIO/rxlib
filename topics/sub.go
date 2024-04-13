package topics

import (
	"fmt"
	"strings"
)

type SubscribeTopicBuilder struct {
	parts []string
}

func NewSubscribeTopicBuilder() *SubscribeTopicBuilder {
	return &SubscribeTopicBuilder{
		parts: []string{},
	}
}

func (b *SubscribeTopicBuilder) WithVersion(version string) *SubscribeTopicBuilder {
	b.parts = append(b.parts, version)
	return b
}

func (b *SubscribeTopicBuilder) WithClientType(clientType ClientType) *SubscribeTopicBuilder {
	b.parts = append(b.parts, string(clientType))
	return b
}

func (b *SubscribeTopicBuilder) WithTargetUUID(targetUUID string) *SubscribeTopicBuilder {
	b.parts = append(b.parts, targetUUID)
	return b
}

func (b *SubscribeTopicBuilder) WithSenderUUID(senderUUID string) *SubscribeTopicBuilder {
	b.parts = append(b.parts, senderUUID)
	return b
}

func (b *SubscribeTopicBuilder) WithRequestUUID(requestUUID string) *SubscribeTopicBuilder {
	b.parts = append(b.parts, requestUUID)
	return b
}

func (b *SubscribeTopicBuilder) WithDataType(dataType DataType) *SubscribeTopicBuilder {
	if dataType == "" {
		dataType = Plain
	}
	b.parts = append(b.parts, string(dataType))
	return b
}

func (b *SubscribeTopicBuilder) WithType(t Type) *SubscribeTopicBuilder {
	b.parts = append(b.parts, string(t))
	return b
}

func (b *SubscribeTopicBuilder) Build() string {
	return fmt.Sprintf("r/resp/%s", strings.Join(b.parts, "/"))
}
