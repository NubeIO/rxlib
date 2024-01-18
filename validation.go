package rxlib

import (
	"fmt"
	"time"
)

type TypeValidation string

const (
	TypeHalt        TypeValidation = "halt"
	TypeError       TypeValidation = "error"
	TypeValidations TypeValidation = "validation"
)

type ErrorsAndValidation struct {
	Type        TypeValidation `json:"type"`
	ObjectError *NewValidation `json:"error"`
	Halt        *NewValidation `json:"haltReason,omitempty"`
	Validation  *NewValidation `json:"validationMessage,omitempty"`
	Custom      *NewValidation
}

type NewValidation struct {
	Message     string    `json:"message"`
	Explanation string    `json:"explanation"`
	Error       error     `json:"-"`
	Timestamp   time.Time `json:"timestamp"`
	Timesince   string    `json:"timesince,omitempty"`
}

type ValidationBuilder interface {
	CustomValidation(v TypeValidation, m *ValidationMessage) ValidationBuilder
	AddValidation(key string) ValidationBuilder
	RemoveValidation(key string) ValidationBuilder
	GetValidations() map[string]*ErrorsAndValidation
	GetValidation(key string) (*ErrorsAndValidation, bool)
}

type ObjectBuilder struct {
	validations map[string]*ErrorsAndValidation
}

func NewObjectBuilder() ValidationBuilder {
	return &ObjectBuilder{
		validations: make(map[string]*ErrorsAndValidation),
	}
}

func (builder *ObjectBuilder) AddValidation(key string) ValidationBuilder {
	validation := &ErrorsAndValidation{}
	builder.validations[key] = validation
	return builder
}

func (builder *ObjectBuilder) CustomValidation(v TypeValidation, m *ValidationMessage) ValidationBuilder {
	builder.validations[string(v)] = &ErrorsAndValidation{
		Type:   v,
		Custom: &NewValidation{Message: m.Message, Explanation: m.Explanation, Error: m.Error, Timestamp: time.Now().UTC()},
	}

	return builder
}

func (builder *ObjectBuilder) RemoveValidation(key string) ValidationBuilder {
	delete(builder.validations, key)
	return builder
}

func (builder *ObjectBuilder) GetValidationWithTimeSince(key string) (*ErrorsAndValidation, bool) {
	v, exists := builder.validations[key]
	if exists {
		if v.ObjectError != nil {
			v.ObjectError.Timesince = TimeSince(v.ObjectError.Timestamp)
		}
		if v.Halt != nil {
			v.Halt.Timesince = TimeSince(v.Halt.Timestamp)
		}
		if v.Validation != nil {
			v.Validation.Timesince = TimeSince(v.Validation.Timestamp)
		}
		return v, true
	}
	return nil, false

}

func (builder *ObjectBuilder) GetValidationsWithTimeSince() map[string]*ErrorsAndValidation {
	for key, v := range builder.validations {
		if v.ObjectError != nil {
			v.ObjectError.Timesince = TimeSince(v.ObjectError.Timestamp)
		}
		if v.Halt != nil {
			v.Halt.Timesince = TimeSince(v.Halt.Timestamp)
		}
		if v.Validation != nil {
			v.Validation.Timesince = TimeSince(v.Validation.Timestamp)
		}
		builder.validations[key] = v
	}
	return builder.validations
}

func (builder *ObjectBuilder) GetValidations() map[string]*ErrorsAndValidation {
	return builder.validations
}

func (builder *ObjectBuilder) GetValidation(key string) (*ErrorsAndValidation, bool) {
	v, exists := builder.validations[key]
	if exists {
		return v, true
	}
	return nil, false
}

type ValidationMessage struct {
	Error       error  `json:"error"`
	Message     string `json:"message"`
	Explanation string `json:"explanation"`
}

func (ev *ErrorsAndValidation) SetError(m *ValidationMessage) {
	ev.Type = TypeError
	var message string
	if m.Error != nil {
		message = m.Error.Error()
	} else {
		message = m.Message
	}
	ev.ObjectError = &NewValidation{
		Error:       m.Error,
		Message:     message,
		Explanation: m.Explanation,
		Timestamp:   time.Now().UTC(),
	}
}

func (ev *ErrorsAndValidation) SetHaltReason(m *ValidationMessage) {
	ev.Type = TypeHalt
	var message string
	if m.Error != nil {
		message = m.Error.Error()
	} else {
		message = m.Message
	}
	ev.Halt = &NewValidation{
		Error:       m.Error,
		Message:     message,
		Explanation: m.Explanation,
		Timestamp:   time.Now().UTC(),
	}
}

func (ev *ErrorsAndValidation) SetValidation(m *ValidationMessage) {
	ev.Type = TypeValidations
	var message string
	if m.Error != nil {
		message = m.Error.Error()
	} else {
		message = m.Message
	}
	ev.Validation = &NewValidation{
		Error:       m.Error,
		Message:     message,
		Explanation: m.Explanation,
		Timestamp:   time.Now().UTC(),
	}
}

func (ev *ErrorsAndValidation) ToString() string {
	return fmt.Sprintf("&{Type:%s Error:%v Halt:%v  Validation:%v}",
		ev.Type,
		ev.ObjectError,
		ev.Halt,
		ev.Validation)
}
