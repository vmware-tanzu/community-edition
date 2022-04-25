// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// AWSAvailabilityZone a w s availability zone
// swagger:model AWSAvailabilityZone
type AWSAvailabilityZone struct {

	// id
	ID string `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this a w s availability zone
func (m *AWSAvailabilityZone) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AWSAvailabilityZone) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AWSAvailabilityZone) UnmarshalBinary(b []byte) error {
	var res AWSAvailabilityZone
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
