// Code generated by go-swagger; DO NOT EDIT.

package unmanaged

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteUnmanagedClusterParams creates a new DeleteUnmanagedClusterParams object
// no default values defined in spec.
func NewDeleteUnmanagedClusterParams() DeleteUnmanagedClusterParams {

	return DeleteUnmanagedClusterParams{}
}

// DeleteUnmanagedClusterParams contains all the bound params for the delete unmanaged cluster operation
// typically these are obtained from a http.Request
//
// swagger:parameters deleteUnmanagedCluster
type DeleteUnmanagedClusterParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The name of the unmanaged cluster.
	  Required: true
	  In: path
	*/
	ClusterName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteUnmanagedClusterParams() beforehand.
func (o *DeleteUnmanagedClusterParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rClusterName, rhkClusterName, _ := route.Params.GetOK("clusterName")
	if err := o.bindClusterName(rClusterName, rhkClusterName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindClusterName binds and validates parameter ClusterName from path.
func (o *DeleteUnmanagedClusterParams) bindClusterName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ClusterName = raw

	return nil
}
