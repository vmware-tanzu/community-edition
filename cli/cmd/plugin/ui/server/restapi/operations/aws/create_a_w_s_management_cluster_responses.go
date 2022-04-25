// Code generated by go-swagger; DO NOT EDIT.

package aws

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// CreateAWSManagementClusterOKCode is the HTTP code returned for type CreateAWSManagementClusterOK
const CreateAWSManagementClusterOKCode int = 200

/*CreateAWSManagementClusterOK Creating management cluster started successfully

swagger:response createAWSManagementClusterOK
*/
type CreateAWSManagementClusterOK struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewCreateAWSManagementClusterOK creates CreateAWSManagementClusterOK with default headers values
func NewCreateAWSManagementClusterOK() *CreateAWSManagementClusterOK {

	return &CreateAWSManagementClusterOK{}
}

// WithPayload adds the payload to the create a w s management cluster o k response
func (o *CreateAWSManagementClusterOK) WithPayload(payload string) *CreateAWSManagementClusterOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create a w s management cluster o k response
func (o *CreateAWSManagementClusterOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAWSManagementClusterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// CreateAWSManagementClusterBadRequestCode is the HTTP code returned for type CreateAWSManagementClusterBadRequest
const CreateAWSManagementClusterBadRequestCode int = 400

/*CreateAWSManagementClusterBadRequest Bad request

swagger:response createAWSManagementClusterBadRequest
*/
type CreateAWSManagementClusterBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateAWSManagementClusterBadRequest creates CreateAWSManagementClusterBadRequest with default headers values
func NewCreateAWSManagementClusterBadRequest() *CreateAWSManagementClusterBadRequest {

	return &CreateAWSManagementClusterBadRequest{}
}

// WithPayload adds the payload to the create a w s management cluster bad request response
func (o *CreateAWSManagementClusterBadRequest) WithPayload(payload *models.Error) *CreateAWSManagementClusterBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create a w s management cluster bad request response
func (o *CreateAWSManagementClusterBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAWSManagementClusterBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateAWSManagementClusterUnauthorizedCode is the HTTP code returned for type CreateAWSManagementClusterUnauthorized
const CreateAWSManagementClusterUnauthorizedCode int = 401

/*CreateAWSManagementClusterUnauthorized Incorrect credentials

swagger:response createAWSManagementClusterUnauthorized
*/
type CreateAWSManagementClusterUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateAWSManagementClusterUnauthorized creates CreateAWSManagementClusterUnauthorized with default headers values
func NewCreateAWSManagementClusterUnauthorized() *CreateAWSManagementClusterUnauthorized {

	return &CreateAWSManagementClusterUnauthorized{}
}

// WithPayload adds the payload to the create a w s management cluster unauthorized response
func (o *CreateAWSManagementClusterUnauthorized) WithPayload(payload *models.Error) *CreateAWSManagementClusterUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create a w s management cluster unauthorized response
func (o *CreateAWSManagementClusterUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAWSManagementClusterUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateAWSManagementClusterInternalServerErrorCode is the HTTP code returned for type CreateAWSManagementClusterInternalServerError
const CreateAWSManagementClusterInternalServerErrorCode int = 500

/*CreateAWSManagementClusterInternalServerError Internal server error

swagger:response createAWSManagementClusterInternalServerError
*/
type CreateAWSManagementClusterInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateAWSManagementClusterInternalServerError creates CreateAWSManagementClusterInternalServerError with default headers values
func NewCreateAWSManagementClusterInternalServerError() *CreateAWSManagementClusterInternalServerError {

	return &CreateAWSManagementClusterInternalServerError{}
}

// WithPayload adds the payload to the create a w s management cluster internal server error response
func (o *CreateAWSManagementClusterInternalServerError) WithPayload(payload *models.Error) *CreateAWSManagementClusterInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create a w s management cluster internal server error response
func (o *CreateAWSManagementClusterInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateAWSManagementClusterInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
