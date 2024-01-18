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
	Type              ErrorsValidation `json:"type"`
	Error             *NewValidation   `json:"-"`
	ErrorMessage      *NewValidation   `json:"errorMessage,omitempty"`
	HaltReason        *NewValidation   `json:"haltReason,omitempty"`
	HaltExplanation   *NewValidation   `json:"haltExplanation,omitempty"`
	ValidationMessage *NewValidation   `json:"validationMessage,omitempty"`
}

type NewValidation struct {
	Message   string    `json:"message"`
	Error     error     `json:"-"`
	Timestamp time.Time `json:"timestamp"`
	Timesince string    `json:"json"`
}

type ValidationBuilder interface {
	AddValidation(key string) ValidationBuilder
	RemoveValidation(key string) *ObjectBuilder

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

func (ev *ErrorsAndValidation) SetError(err error) {
	ev.Type = TypeError
	ev.Error = &NewValidation{
		Error:     err,
		Timestamp: time.Now().UTC(),
	}
}

func (builder *ObjectBuilder) RemoveValidation(key string) *ObjectBuilder {
	delete(builder.validations, key)
	return builder
}

func (builder *ObjectBuilder) GetValidationWithTimeSince(key string) (*ErrorsAndValidation, bool) {
	v, exists := builder.validations[key]
	if exists {
		if v.Error != nil {
			v.Error.Timesince = timeSince(v.Error.Timestamp)
		}
		if v.ErrorMessage != nil {
			v.ErrorMessage.Timesince = timeSince(v.ErrorMessage.Timestamp)
		}
		if v.HaltReason != nil {
			v.HaltReason.Timesince = timeSince(v.HaltReason.Timestamp)
		}
		if v.HaltExplanation != nil {
			v.HaltExplanation.Timesince = timeSince(v.HaltExplanation.Timestamp)
		}
		if v.ValidationMessage != nil {
			v.ValidationMessage.Timesince = timeSince(v.ValidationMessage.Timestamp)
		}
		return v, true
	}
	return nil, false

}

func (builder *ObjectBuilder) GetValidationsWithTimeSince() map[string]*ErrorsAndValidation {
	for key, v := range builder.validations {
		if v.Error != nil {
			v.Error.Timesince = timeSince(v.Error.Timestamp)
		}
		if v.ErrorMessage != nil {
			v.ErrorMessage.Timesince = timeSince(v.ErrorMessage.Timestamp)
		}
		if v.HaltReason != nil {
			v.HaltReason.Timesince = timeSince(v.HaltReason.Timestamp)
		}
		if v.HaltExplanation != nil {
			v.HaltExplanation.Timesince = timeSince(v.HaltExplanation.Timestamp)
		}
		if v.ValidationMessage != nil {
			v.ValidationMessage.Timesince = timeSince(v.ValidationMessage.Timestamp)
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

func (ev *ErrorsAndValidation) SetErrorMessage(message string) {
	ev.ErrorMessage = &NewValidation{
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
}

func (ev *ErrorsAndValidation) SetHaltReason(reason string) {
	ev.HaltReason = &NewValidation{
		Message:   reason,
		Timestamp: time.Now().UTC(),
	}
}

func (ev *ErrorsAndValidation) SetHaltExplanation(explanation string) {
	ev.HaltExplanation = &NewValidation{
		Message:   explanation,
		Timestamp: time.Now().UTC(),
	}
}

func (ev *ErrorsAndValidation) SetValidationMessage(message string) {
	ev.ValidationMessage = &NewValidation{
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
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
