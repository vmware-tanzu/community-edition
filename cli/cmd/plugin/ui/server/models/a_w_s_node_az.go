// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AWSNodeAz a w s node az
//
// swagger:model AWSNodeAz
type AWSNodeAz struct {

	// name
	Name string `json:"name,omitempty"`

	// private subnet ID
	PrivateSubnetID string `json:"privateSubnetID,omitempty"`

	// public subnet ID
	PublicSubnetID string `json:"publicSubnetID,omitempty"`

	// worker node type
	WorkerNodeType string `json:"workerNodeType,omitempty"`
}

// Validate validates this a w s node az
func (m *AWSNodeAz) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this a w s node az based on context it is used
func (m *AWSNodeAz) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AWSNodeAz) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AWSNodeAz) UnmarshalBinary(b []byte) error {
	var res AWSNodeAz
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}