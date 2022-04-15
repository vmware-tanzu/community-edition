// Code generated by go-swagger; DO NOT EDIT.

package vsphere

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// GetVsphereThumbprintOKCode is the HTTP code returned for type GetVsphereThumbprintOK
const GetVsphereThumbprintOKCode int = 200

/*GetVsphereThumbprintOK Successful retrieval of vSphere thumbprint

swagger:response getVsphereThumbprintOK
*/
type GetVsphereThumbprintOK struct {

	/*
	  In: Body
	*/
	Payload *models.VSphereThumbprint `json:"body,omitempty"`
}

// NewGetVsphereThumbprintOK creates GetVsphereThumbprintOK with default headers values
func NewGetVsphereThumbprintOK() *GetVsphereThumbprintOK {

	return &GetVsphereThumbprintOK{}
}

// WithPayload adds the payload to the get vsphere thumbprint o k response
func (o *GetVsphereThumbprintOK) WithPayload(payload *models.VSphereThumbprint) *GetVsphereThumbprintOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get vsphere thumbprint o k response
func (o *GetVsphereThumbprintOK) SetPayload(payload *models.VSphereThumbprint) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVsphereThumbprintOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetVsphereThumbprintBadRequestCode is the HTTP code returned for type GetVsphereThumbprintBadRequest
const GetVsphereThumbprintBadRequestCode int = 400

/*GetVsphereThumbprintBadRequest Bad request

swagger:response getVsphereThumbprintBadRequest
*/
type GetVsphereThumbprintBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVsphereThumbprintBadRequest creates GetVsphereThumbprintBadRequest with default headers values
func NewGetVsphereThumbprintBadRequest() *GetVsphereThumbprintBadRequest {

	return &GetVsphereThumbprintBadRequest{}
}

// WithPayload adds the payload to the get vsphere thumbprint bad request response
func (o *GetVsphereThumbprintBadRequest) WithPayload(payload *models.Error) *GetVsphereThumbprintBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get vsphere thumbprint bad request response
func (o *GetVsphereThumbprintBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVsphereThumbprintBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetVsphereThumbprintUnauthorizedCode is the HTTP code returned for type GetVsphereThumbprintUnauthorized
const GetVsphereThumbprintUnauthorizedCode int = 401

/*GetVsphereThumbprintUnauthorized Incorrect credentials

swagger:response getVsphereThumbprintUnauthorized
*/
type GetVsphereThumbprintUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVsphereThumbprintUnauthorized creates GetVsphereThumbprintUnauthorized with default headers values
func NewGetVsphereThumbprintUnauthorized() *GetVsphereThumbprintUnauthorized {

	return &GetVsphereThumbprintUnauthorized{}
}

// WithPayload adds the payload to the get vsphere thumbprint unauthorized response
func (o *GetVsphereThumbprintUnauthorized) WithPayload(payload *models.Error) *GetVsphereThumbprintUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get vsphere thumbprint unauthorized response
func (o *GetVsphereThumbprintUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVsphereThumbprintUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetVsphereThumbprintInternalServerErrorCode is the HTTP code returned for type GetVsphereThumbprintInternalServerError
const GetVsphereThumbprintInternalServerErrorCode int = 500

/*GetVsphereThumbprintInternalServerError Internal server error

swagger:response getVsphereThumbprintInternalServerError
*/
type GetVsphereThumbprintInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetVsphereThumbprintInternalServerError creates GetVsphereThumbprintInternalServerError with default headers values
func NewGetVsphereThumbprintInternalServerError() *GetVsphereThumbprintInternalServerError {

	return &GetVsphereThumbprintInternalServerError{}
}

// WithPayload adds the payload to the get vsphere thumbprint internal server error response
func (o *GetVsphereThumbprintInternalServerError) WithPayload(payload *models.Error) *GetVsphereThumbprintInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get vsphere thumbprint internal server error response
func (o *GetVsphereThumbprintInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetVsphereThumbprintInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
