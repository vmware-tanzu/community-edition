package errors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		errorFunc func(err error, args ...interface{}) TanzuError
		Reason    Reason
	}{
		{
			NewLockFailed, ReasonLockFailed,
		},
		{
			NewMergeFailed, ReasonMergeFailed,
		},
		{
			NewReadFailed, ReasonReadFailed,
		},
		{
			NewIOFailed, ReasonIOFailed,
		},
		{
			NewMarshallingFailed, ReasonMarshallingFailed,
		},
		{
			NewNormalizationFailed, ReasonNormalizationFailed,
		},
		{
			NewValidationFailed, ReasonValidationFailed,
		},
		{
			NewTimeout, ReasonTimeout,
		},
		{
			NewGenericFail, ReasonGeneric,
		},
	}

	for _, tc := range tests {
		tc := tc
		name := runtime.FuncForPC(reflect.ValueOf(tc.errorFunc).Pointer()).Name()
		t.Run(name, func(t *testing.T) {
			err := tc.errorFunc(errors.New("fail"), "some failure. Arg %s", "test")
			if err.ErrorDetails.Reason != tc.Reason {
				t.Errorf("%s, Expected: %v, got: %v", name, tc.Reason, err.ErrorDetails.Reason)
			}
			expected := fmt.Sprintf("%s: %s (%s)", tc.Reason, "some failure. Arg test", "fail")
			if err.Error() != expected {
				t.Errorf("%s, Expected: %s, got: %s", name, expected, err.Error())
			}
		})
	}
}
