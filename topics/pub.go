package topics

import (
	"fmt"
	"strings"
)

type DataType string

const (
	Plain DataType = "plain"
	Proto DataType = "proto"
)

type ClientType string

const (
	Edge  ClientType = "edge"
	Cloud ClientType = "cloud"
	UI    ClientType = "ui"
)

type Type string

const (
	Command Type = "command"
	Object  Type = "object"
	Objects Type = "objects"
	Runtime Type = "runtime"
)

type PublishTopicBuilder struct {
	parts []string
}

func NewPublishTopicBuilder() *PublishTopicBuilder {
	return &PublishTopicBuilder{
		parts: []string{},
	}
}

func (b *PublishTopicBuilder) WithVersion(version string) *PublishTopicBuilder {
	b.parts = append(b.parts, version)
	return b
}

func (b *PublishTopicBuilder) WithClientType(clientType ClientType) *PublishTopicBuilder {
	b.parts = append(b.parts, string(clientType))
	return b
}

func (b *PublishTopicBuilder) WithTargetUUID(targetUUID string) *PublishTopicBuilder {
	b.parts = append(b.parts, targetUUID)
	return b
}

func (b *PublishTopicBuilder) WithSenderUUID(senderUUID string) *PublishTopicBuilder {
	b.parts = append(b.parts, senderUUID)
	return b
}

func (b *PublishTopicBuilder) WithRequestUUID(requestUUID string) *PublishTopicBuilder {
	b.parts = append(b.parts, requestUUID)
	return b
}

func (b *PublishTopicBuilder) WithDataType(dataType DataType) *PublishTopicBuilder {
	if dataType == "" {
		dataType = Plain
	}
	b.parts = append(b.parts, string(dataType))
	return b
}

func (b *PublishTopicBuilder) WithType(t Type) *PublishTopicBuilder {
	b.parts = append(b.parts, string(t))
	return b
}

func (b *PublishTopicBuilder) Build() string {
	return fmt.Sprintf("r/req/%s", strings.Join(b.parts, "/"))
}
