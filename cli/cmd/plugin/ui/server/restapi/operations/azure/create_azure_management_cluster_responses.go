// Code generated by go-swagger; DO NOT EDIT.

package azure

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// CreateAzureManagementClusterOKCode is the HTTP code returned for type CreateAzureManagementClusterOK
const CreateAzureManagementClusterOKCode int = 200

/*CreateAzureManagementClusterOK Creating management cluster started successfully

swagger:response createAzureManagementClusterOK
*/
type CreateAzureManagementClusterOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateAzureManagementClusterOK creates CreateAzureManagementClusterOK with default headers values
func NewCreateAzureManagementClusterOK() *CreateAzureManagementClusterOK {

	return &CreateAzureManagementClusterOK{}
}

// WithPayload adds the payload to the create azure management cluster o k response
func (o *CreateAzureManagementClusterOK) WithPayload(payload string) *CreateAzureManagementClusterOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create azure management cluster o k response
func (o *CreateAzureManagementClusterOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAzureManagementClusterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// CreateAzureManagementClusterBadRequestCode is the HTTP code returned for type CreateAzureManagementClusterBadRequest
const CreateAzureManagementClusterBadRequestCode int = 400

/*CreateAzureManagementClusterBadRequest Bad request

swagger:response createAzureManagementClusterBadRequest
*/
type CreateAzureManagementClusterBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateAzureManagementClusterBadRequest creates CreateAzureManagementClusterBadRequest with default headers values
func NewCreateAzureManagementClusterBadRequest() *CreateAzureManagementClusterBadRequest {

	return &CreateAzureManagementClusterBadRequest{}
}

// WithPayload adds the payload to the create azure management cluster bad request response
func (o *CreateAzureManagementClusterBadRequest) WithPayload(payload *models.Error) *CreateAzureManagementClusterBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create azure management cluster bad request response
func (o *CreateAzureManagementClusterBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAzureManagementClusterBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateAzureManagementClusterUnauthorizedCode is the HTTP code returned for type CreateAzureManagementClusterUnauthorized
const CreateAzureManagementClusterUnauthorizedCode int = 401

/*CreateAzureManagementClusterUnauthorized Incorrect credentials

swagger:response createAzureManagementClusterUnauthorized
*/
type CreateAzureManagementClusterUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateAzureManagementClusterUnauthorized creates CreateAzureManagementClusterUnauthorized with default headers values
func NewCreateAzureManagementClusterUnauthorized() *CreateAzureManagementClusterUnauthorized {

	return &CreateAzureManagementClusterUnauthorized{}
}

// WithPayload adds the payload to the create azure management cluster unauthorized response
func (o *CreateAzureManagementClusterUnauthorized) WithPayload(payload *models.Error) *CreateAzureManagementClusterUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create azure management cluster unauthorized response
func (o *CreateAzureManagementClusterUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAzureManagementClusterUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateAzureManagementClusterInternalServerErrorCode is the HTTP code returned for type CreateAzureManagementClusterInternalServerError
const CreateAzureManagementClusterInternalServerErrorCode int = 500

/*CreateAzureManagementClusterInternalServerError Internal server error

swagger:response createAzureManagementClusterInternalServerError
*/
type CreateAzureManagementClusterInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateAzureManagementClusterInternalServerError creates CreateAzureManagementClusterInternalServerError with default headers values
func NewCreateAzureManagementClusterInternalServerError() *CreateAzureManagementClusterInternalServerError {

	return &CreateAzureManagementClusterInternalServerError{}
}

// WithPayload adds the payload to the create azure management cluster internal server error response
func (o *CreateAzureManagementClusterInternalServerError) WithPayload(payload *models.Error) *CreateAzureManagementClusterInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create azure management cluster internal server error response
func (o *CreateAzureManagementClusterInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAzureManagementClusterInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
