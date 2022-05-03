// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// WorkloadCluster workload cluster
// swagger:model WorkloadCluster
type WorkloadCluster struct {

	// cpcount
	Cpcount string `json:"cpcount,omitempty"`

	// k8sversion
	K8sversion string `json:"k8sversion,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// namespace
	Namespace string `json:"namespace,omitempty"`

	// plan
	Plan string `json:"plan,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// wncount
	Wncount string `json:"wncount,omitempty"`
}

// Validate validates this workload cluster
func (m *WorkloadCluster) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *WorkloadCluster) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WorkloadCluster) UnmarshalBinary(b []byte) error {
	var res WorkloadCluster
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
