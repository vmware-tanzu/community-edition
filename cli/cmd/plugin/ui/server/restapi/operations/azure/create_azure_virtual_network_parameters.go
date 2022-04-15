// Code generated by go-swagger; DO NOT EDIT.

package azure

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

// NewCreateAzureVirtualNetworkParams creates a new CreateAzureVirtualNetworkParams object
// no default values defined in spec.
func NewCreateAzureVirtualNetworkParams() CreateAzureVirtualNetworkParams {

	return CreateAzureVirtualNetworkParams{}
}

// CreateAzureVirtualNetworkParams contains all the bound params for the create azure virtual network operation
// typically these are obtained from a http.Request
//
// swagger:parameters createAzureVirtualNetwork
type CreateAzureVirtualNetworkParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*parameters to create a new Azure Virtual network
	  Required: true
	  In: body
	*/
	Params *models.AzureVirtualNetwork
	/*Name of the Azure resource group
	  Required: true
	  In: path
	*/
	ResourceGroupName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateAzureVirtualNetworkParams() beforehand.
func (o *CreateAzureVirtualNetworkParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.AzureVirtualNetwork
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
	rResourceGroupName, rhkResourceGroupName, _ := route.Params.GetOK("resourceGroupName")
	if err := o.bindResourceGroupName(rResourceGroupName, rhkResourceGroupName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindResourceGroupName binds and validates parameter ResourceGroupName from path.
func (o *CreateAzureVirtualNetworkParams) bindResourceGroupName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ResourceGroupName = raw

	return nil
}
