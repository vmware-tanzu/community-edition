package errors

import (
	"fmt"
)

type Reason string

const (
	ReasonUnknown             Reason = "Unknown"
	ReasonLockFailed          Reason = "LockFailed"
	ReasonMergeFailed         Reason = "MergeFailed"
	ReasonReadFailed          Reason = "ReadFailed"
	ReasonIOFailed            Reason = "IOFailed"
	ReasonMarshallingFailed   Reason = "MarshallingFailed"
	ReasonEncodeFailed        Reason = "EncodeFailed"
	ReasonNormalizationFailed Reason = "NormalizationFailed"
	ReasonValidationFailed    Reason = "ValidationFailed"
	ReasonTimeout             Reason = "Timeout"
	ReasonGeneric             Reason = "GenericFailure"
)

// TODO: implement the errors.Unwrap interface to support errors.Unwrap cleanly
type TanzuError struct {
	ErrorDetails ErrorDetails `json:"errorDetails"`
}

type ErrorDetails struct {
	Message string `json:"message,omitempty"`
	Reason  Reason `json:"reason,omitempty"`
	Err     error  `json:"err"`
}

func (t TanzuError) Error() string {
	return fmt.Sprintf("%s: %s (%s)", t.ErrorDetails.Reason, t.ErrorDetails.Message, t.ErrorDetails.Err)
}

func generateError(reason Reason, err error, args ...interface{}) TanzuError {
	var message string
	switch len(args) {
	case 0:
		message = ""
	case 1:
		message = args[0].(string)
	default:
		message = fmt.Sprintf(args[0].(string), args[1:]...)
	}
	return TanzuError{
		ErrorDetails: ErrorDetails{
			Reason:  reason,
			Message: message,
			Err:     err,
		},
	}
}

func NewLockFailed(err error, args ...interface{}) TanzuError {
	return generateError(ReasonLockFailed, err, args...)
}

func NewMergeFailed(err error, args ...interface{}) TanzuError {
	return generateError(ReasonMergeFailed, err, args...)
}

func NewReadFailed(err error, args ...interface{}) TanzuError {
	return generateError(ReasonReadFailed, err, args...)
}

func NewIOFailed(err error, args ...interface{}) TanzuError {
	return generateError(ReasonIOFailed, err, args...)
}

func NewMarshallingFailed(err error, args ...interface{}) TanzuError {
	return generateError(ReasonMarshallingFailed, err, args...)
}

func NewNormalizationFailed(err error, args ...interface{}) TanzuError {
	return generateError(ReasonNormalizationFailed, err, args...)
}

func NewValidationFailed(err error, args ...interface{}) TanzuError {
	return generateError(ReasonValidationFailed, err, args...)
}

func NewTimeout(err error, args ...interface{}) TanzuError {
	return generateError(ReasonTimeout, err, args...)
}

func NewGenericFail(err error, args ...interface{}) TanzuError {
	return generateError(ReasonGeneric, err, args...)
}

func NewEncodeFail(err error, args ...interface{}) TanzuError {
	return generateError(ReasonEncodeFailed, err, args...)
}
