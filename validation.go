package rxlib

import (
	"time"
)

type TypeValidation string

const (
	TypeHalt        TypeValidation = "halt"
	TypeError       TypeValidation = "Err"
	TypeValidations TypeValidation = "validation"
)

type ErrorsAndValidation struct {
	Type        TypeValidation `json:"type"`
	ObjectError *NewValidation `json:"Err"`
	Halt        *NewValidation `json:"haltReason,omitempty"`
	Validation  *NewValidation `json:"validationMessage,omitempty"`
	Custom      *NewValidation `json:"custom,omitempty"`
}

type NewValidation struct {
	Message     string    `json:"message"`
	Explanation string    `json:"explanation,omitempty"`
	Error       error     `json:"-"`
	Timestamp   time.Time `json:"timestamp"`
	Timesince   string    `json:"timesince,omitempty"`
}

type ValidationMessage struct {
	Error       error  `json:"Err,omitempty"`
	Message     string `json:"message,omitempty"`
	Explanation string `json:"explanation,omitempty"`
}

//
//type ValidationBuilder interface {
//	CustomValidation(v TypeValidation, m *ValidationMessage) (*ErrorsAndValidation, Err)
//	CreateValidation(m *ValidationMessage) ValidationBuilder
//}
//
//type ObjectBuilder struct {
//}
//
//func NewValidationBuilder() ValidationBuilder {
//	return &ObjectBuilder{}
//}
//
//func (inst *ObjectBuilder) CreateValidation(m *ValidationMessage) ValidationBuilder {
//	ev := &ErrorsAndValidation{}
//	var message string
//	if m.Error != nil {
//		message = m.Error.Error()
//	} else {
//		message = m.Payload
//	}
//	ev.Custom = &NewValidation{
//		Error:       m.Error,
//		Payload:     message,
//		Explanation: m.Explanation,
//		Timestamp:   time.Now().UTC(),
//	}
//	return inst
//}
//

//
//func (inst *ObjectBuilder) CustomValidation(v TypeValidation, m *ValidationMessage) (*ErrorsAndValidation, Err) {
//	if v == "" {
//		return nil, errors.New("TypeValidation can not be empty")
//	}
//	if m == nil {
//		return nil, errors.New("*ValidationMessage can not be nil")
//	}
//	return &ErrorsAndValidation{
//		Type:   v,
//		Custom: &NewValidation{Payload: m.Payload, Explanation: m.Explanation, Error: m.Error, Timestamp: time.Now().UTC()},
//	}, nil
//}
//
//func (ev *ErrorsAndValidation) SetError(m *ValidationMessage) {
//	ev.Type = TypeError
//	var message string
//	if m.Error != nil {
//		message = m.Error.Error()
//	} else {
//		message = m.Payload
//	}
//	ev.ObjectError = &NewValidation{
//		Error:       m.Error,
//		Payload:     message,
//		Explanation: m.Explanation,
//		Timestamp:   time.Now().UTC(),
//	}
//}
//
//func (ev *ErrorsAndValidation) SetHaltReason(m *ValidationMessage) {
//	ev.Type = TypeHalt
//	var message string
//	if m.Error != nil {
//		message = m.Error.Error()
//	} else {
//		message = m.Payload
//	}
//	ev.Halt = &NewValidation{
//		Error:       m.Error,
//		Payload:     message,
//		Explanation: m.Explanation,
//		Timestamp:   time.Now().UTC(),
//	}
//}
//
//func (ev *ErrorsAndValidation) SetValidation(m *ValidationMessage) {
//	ev.Type = TypeValidations
//	var message string
//	if m.Error != nil {
//		message = m.Error.Error()
//	} else {
//		message = m.Payload
//	}
//	ev.Validation = &NewValidation{
//		Error:       m.Error,
//		Payload:     message,
//		Explanation: m.Explanation,
//		Timestamp:   time.Now().UTC(),
//	}
//}
