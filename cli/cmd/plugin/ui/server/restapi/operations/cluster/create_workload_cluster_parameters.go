// Code generated by go-swagger; DO NOT EDIT.

package cluster

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// NewCreateWorkloadClusterParams creates a new CreateWorkloadClusterParams object
// no default values defined in spec.
func NewCreateWorkloadClusterParams() CreateWorkloadClusterParams {

	return CreateWorkloadClusterParams{}
}

// CreateWorkloadClusterParams contains all the bound params for the create workload cluster operation
// typically these are obtained from a http.Request
//
// swagger:parameters createWorkloadCluster
type CreateWorkloadClusterParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The name of the management cluster.
	  Required: true
	  In: path
	*/
	ManagementClusterName string
	/*Parameters to create a workload cluster.
	  Required: true
	  In: body
	*/
	Params *models.CreateWorkloadClusterParams
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateWorkloadClusterParams() beforehand.
func (o *CreateWorkloadClusterParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rManagementClusterName, rhkManagementClusterName, _ := route.Params.GetOK("managementClusterName")
	if err := o.bindManagementClusterName(rManagementClusterName, rhkManagementClusterName, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.CreateWorkloadClusterParams
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("params", "body"))
			} else {
				res = append(res, errors.NewParseError("params", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Params = &body
			}
		}
	} else {
		res = append(res, errors.Required("params", "body"))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindManagementClusterName binds and validates parameter ManagementClusterName from path.
func (o *CreateWorkloadClusterParams) bindManagementClusterName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ManagementClusterName = raw

	return nil
}