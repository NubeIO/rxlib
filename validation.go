package rxlib

import (
	"fmt"
	"time"
)

type ErrorsValidation string

const (
	TypeHalt       ErrorsValidation = "halt"
	TypeError      ErrorsValidation = "error"
	TypeValidation ErrorsValidation = "validation"
)

type ErrorsAndValidation struct {
	Type        ErrorsValidation `json:"type"`
	ObjectError *NewValidation   `json:"error"`
	Halt        *NewValidation   `json:"haltReason,omitempty"`
	Validation  *NewValidation   `json:"validationMessage,omitempty"`
}

type NewValidation struct {
	Message     string    `json:"message"`
	Explanation string    `json:"explanation"`
	Error       error     `json:"-"`
	Timestamp   time.Time `json:"timestamp"`
	Timesince   string    `json:"timesince,omitempty"`
}

type ValidationBuilder interface {
	AddValidation(key string) ValidationBuilder
	RemoveValidation(key string) ValidationBuilder
	GetValidations() map[string]*ErrorsAndValidation
	GetValidation(key string) (*ErrorsAndValidation, bool)
}

type ObjectBuilder struct {
	validations map[string]*ErrorsAndValidation
}

func (builder *ObjectBuilder) AddValidation(key string) ValidationBuilder {
	validation := &ErrorsAndValidation{}
	builder.validations[key] = validation
	return builder
}

func NewObjectBuilder() ValidationBuilder {
	return &ObjectBuilder{
		validations: make(map[string]*ErrorsAndValidation),
	}
}

func (builder *ObjectBuilder) RemoveValidation(key string) ValidationBuilder {
	delete(builder.validations, key)
	return builder
}

func (builder *ObjectBuilder) GetValidationWithTimeSince(key string) (*ErrorsAndValidation, bool) {
	v, exists := builder.validations[key]
	if exists {
		if v.ObjectError != nil {
			v.ObjectError.Timesince = timeSince(v.ObjectError.Timestamp)
		}
		if v.Halt != nil {
			v.Halt.Timesince = timeSince(v.Halt.Timestamp)
		}
		if v.Validation != nil {
			v.Validation.Timesince = timeSince(v.Validation.Timestamp)
		}
		return v, true
	}
	return nil, false

}

func (builder *ObjectBuilder) GetValidationsWithTimeSince() map[string]*ErrorsAndValidation {
	for key, v := range builder.validations {
		if v.ObjectError != nil {
			v.ObjectError.Timesince = timeSince(v.ObjectError.Timestamp)
		}
		if v.Halt != nil {
			v.Halt.Timesince = timeSince(v.Halt.Timestamp)
		}
		if v.Validation != nil {
			v.Validation.Timesince = timeSince(v.Validation.Timestamp)
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
	ev.Type = TypeValidation
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

func timeSince(t time.Time) string {
	var duration = time.Since(t)
	switch {
	case duration < 30*time.Second:
		return "just now"
	case duration < 1*time.Minute:
		return fmt.Sprintf("%d sec", int(duration.Seconds()))
	case duration < 1*time.Hour:
		return fmt.Sprintf("%d min ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration < 30*24*time.Hour: // Approximating a month as 30 days
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	case duration < 365*24*time.Hour: // Approximating a year as 365 days
		return fmt.Sprintf("%d months ago", int(duration.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d years ago", int(duration.Hours()/(24*365)))
	}
}
