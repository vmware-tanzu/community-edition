// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// VsphereRegionalClusterParams vsphere regional cluster params
//
// swagger:model VsphereRegionalClusterParams
type VsphereRegionalClusterParams struct {

	// annotations
	Annotations map[string]string `json:"annotations,omitempty"`

	// avi config
	AviConfig *AviConfig `json:"aviConfig,omitempty"`

	// ceip opt in
	CeipOptIn *bool `json:"ceipOptIn,omitempty"`

	// cluster name
	ClusterName string `json:"clusterName,omitempty"`

	// control plane endpoint
	ControlPlaneEndpoint string `json:"controlPlaneEndpoint,omitempty"`

	// control plane flavor
	ControlPlaneFlavor string `json:"controlPlaneFlavor,omitempty"`

	// control plane node type
	ControlPlaneNodeType string `json:"controlPlaneNodeType,omitempty"`

	// datacenter
	Datacenter string `json:"datacenter,omitempty"`

	// datastore
	Datastore string `json:"datastore,omitempty"`

	// enable audit logging
	EnableAuditLogging bool `json:"enableAuditLogging,omitempty"`

	// folder
	Folder string `json:"folder,omitempty"`

	// identity management
	IdentityManagement *IdentityManagementConfig `json:"identityManagement,omitempty"`

	// ip family
	IPFamily string `json:"ipFamily,omitempty"`

	// kubernetes version
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// labels
	Labels map[string]string `json:"labels,omitempty"`

	// machine health check enabled
	MachineHealthCheckEnabled bool `json:"machineHealthCheckEnabled,omitempty"`

	// networking
	Networking *TKGNetwork `json:"networking,omitempty"`

	// num of worker node
	NumOfWorkerNode int64 `json:"numOfWorkerNode,omitempty"`

	// os
	Os *VSphereVirtualMachine `json:"os,omitempty"`

	// resource pool
	ResourcePool string `json:"resourcePool,omitempty"`

	// ssh key
	SSHKey string `json:"ssh_key,omitempty"`

	// vsphere credentials
	VsphereCredentials *VSphereCredentials `json:"vsphereCredentials,omitempty"`

	// worker node type
	WorkerNodeType string `json:"workerNodeType,omitempty"`
}

// Validate validates this vsphere regional cluster params
func (m *VsphereRegionalClusterParams) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAviConfig(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIdentityManagement(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworking(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVsphereCredentials(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VsphereRegionalClusterParams) validateAviConfig(formats strfmt.Registry) error {
	if swag.IsZero(m.AviConfig) { // not required
		return nil
	}

	if m.AviConfig != nil {
		if err := m.AviConfig.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("aviConfig")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("aviConfig")
			}
			return err
		}
	}

	return nil
}

func (m *VsphereRegionalClusterParams) validateIdentityManagement(formats strfmt.Registry) error {
	if swag.IsZero(m.IdentityManagement) { // not required
		return nil
	}

	if m.IdentityManagement != nil {
		if err := m.IdentityManagement.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("identityManagement")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("identityManagement")
			}
			return err
		}
	}

	return nil
}

func (m *VsphereRegionalClusterParams) validateNetworking(formats strfmt.Registry) error {
	if swag.IsZero(m.Networking) { // not required
		return nil
	}

	if m.Networking != nil {
		if err := m.Networking.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("networking")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("networking")
			}
			return err
		}
	}

	return nil
}

func (m *VsphereRegionalClusterParams) validateOs(formats strfmt.Registry) error {
	if swag.IsZero(m.Os) { // not required
		return nil
	}

	if m.Os != nil {
		if err := m.Os.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("os")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("os")
			}
			return err
		}
	}

	return nil
}

func (m *VsphereRegionalClusterParams) validateVsphereCredentials(formats strfmt.Registry) error {
	if swag.IsZero(m.VsphereCredentials) { // not required
		return nil
	}

	if m.VsphereCredentials != nil {
		if err := m.VsphereCredentials.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("vsphereCredentials")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("vsphereCredentials")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this vsphere regional cluster params based on the context it is used
func (m *VsphereRegionalClusterParams) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAviConfig(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIdentityManagement(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNetworking(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateOs(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateVsphereCredentials(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VsphereRegionalClusterParams) contextValidateAviConfig(ctx context.Context, formats strfmt.Registry) error {

	if m.AviConfig != nil {
		if err := m.AviConfig.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("aviConfig")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("aviConfig")
			}
			return err
		}
	}

	return nil
}

func (m *VsphereRegionalClusterParams) contextValidateIdentityManagement(ctx context.Context, formats strfmt.Registry) error {

	if m.IdentityManagement != nil {
		if err := m.IdentityManagement.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("identityManagement")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("identityManagement")
			}
			return err
		}
	}

	return nil
}

func (m *VsphereRegionalClusterParams) contextValidateNetworking(ctx context.Context, formats strfmt.Registry) error {

	if m.Networking != nil {
		if err := m.Networking.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("networking")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("networking")
			}
			return err
		}
	}

	return nil
}

func (m *VsphereRegionalClusterParams) contextValidateOs(ctx context.Context, formats strfmt.Registry) error {

	if m.Os != nil {
		if err := m.Os.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("os")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("os")
			}
			return err
		}
	}

	return nil
}

func (m *VsphereRegionalClusterParams) contextValidateVsphereCredentials(ctx context.Context, formats strfmt.Registry) error {

	if m.VsphereCredentials != nil {
		if err := m.VsphereCredentials.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("vsphereCredentials")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("vsphereCredentials")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VsphereRegionalClusterParams) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VsphereRegionalClusterParams) UnmarshalBinary(b []byte) error {
	var res VsphereRegionalClusterParams
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
