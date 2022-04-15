// Code generated by go-swagger; DO NOT EDIT.

package vsphere

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// SetVSphereEndpointCreatedCode is the HTTP code returned for type SetVSphereEndpointCreated
const SetVSphereEndpointCreatedCode int = 201

/*SetVSphereEndpointCreated Connection successful

swagger:response setVSphereEndpointCreated
*/
type SetVSphereEndpointCreated struct {

	/*
	  In: Body
	*/
	Payload *models.VsphereInfo `json:"body,omitempty"`
}

// NewSetVSphereEndpointCreated creates SetVSphereEndpointCreated with default headers values
func NewSetVSphereEndpointCreated() *SetVSphereEndpointCreated {

	return &SetVSphereEndpointCreated{}
}

// WithPayload adds the payload to the set v sphere endpoint created response
func (o *SetVSphereEndpointCreated) WithPayload(payload *models.VsphereInfo) *SetVSphereEndpointCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set v sphere endpoint created response
func (o *SetVSphereEndpointCreated) SetPayload(payload *models.VsphereInfo) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetVSphereEndpointCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SetVSphereEndpointBadRequestCode is the HTTP code returned for type SetVSphereEndpointBadRequest
const SetVSphereEndpointBadRequestCode int = 400

/*SetVSphereEndpointBadRequest Bad request

swagger:response setVSphereEndpointBadRequest
*/
type SetVSphereEndpointBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSetVSphereEndpointBadRequest creates SetVSphereEndpointBadRequest with default headers values
func NewSetVSphereEndpointBadRequest() *SetVSphereEndpointBadRequest {

	return &SetVSphereEndpointBadRequest{}
}

// WithPayload adds the payload to the set v sphere endpoint bad request response
func (o *SetVSphereEndpointBadRequest) WithPayload(payload *models.Error) *SetVSphereEndpointBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set v sphere endpoint bad request response
func (o *SetVSphereEndpointBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetVSphereEndpointBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SetVSphereEndpointUnauthorizedCode is the HTTP code returned for type SetVSphereEndpointUnauthorized
const SetVSphereEndpointUnauthorizedCode int = 401

/*SetVSphereEndpointUnauthorized Incorrect credentials

swagger:response setVSphereEndpointUnauthorized
*/
type SetVSphereEndpointUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSetVSphereEndpointUnauthorized creates SetVSphereEndpointUnauthorized with default headers values
func NewSetVSphereEndpointUnauthorized() *SetVSphereEndpointUnauthorized {

	return &SetVSphereEndpointUnauthorized{}
}

// WithPayload adds the payload to the set v sphere endpoint unauthorized response
func (o *SetVSphereEndpointUnauthorized) WithPayload(payload *models.Error) *SetVSphereEndpointUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set v sphere endpoint unauthorized response
func (o *SetVSphereEndpointUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetVSphereEndpointUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SetVSphereEndpointInternalServerErrorCode is the HTTP code returned for type SetVSphereEndpointInternalServerError
const SetVSphereEndpointInternalServerErrorCode int = 500

/*SetVSphereEndpointInternalServerError Internal server error

swagger:response setVSphereEndpointInternalServerError
*/
type SetVSphereEndpointInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSetVSphereEndpointInternalServerError creates SetVSphereEndpointInternalServerError with default headers values
func NewSetVSphereEndpointInternalServerError() *SetVSphereEndpointInternalServerError {

	return &SetVSphereEndpointInternalServerError{}
}

// WithPayload adds the payload to the set v sphere endpoint internal server error response
func (o *SetVSphereEndpointInternalServerError) WithPayload(payload *models.Error) *SetVSphereEndpointInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the set v sphere endpoint internal server error response
func (o *SetVSphereEndpointInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SetVSphereEndpointInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
