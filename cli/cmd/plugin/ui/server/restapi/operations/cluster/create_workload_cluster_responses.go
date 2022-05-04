// Code generated by go-swagger; DO NOT EDIT.

package cluster

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// CreateWorkloadClusterOKCode is the HTTP code returned for type CreateWorkloadClusterOK
const CreateWorkloadClusterOKCode int = 200

/*CreateWorkloadClusterOK Create workload cluster started successfully.

swagger:response createWorkloadClusterOK
*/
type CreateWorkloadClusterOK struct {

	/*
	  In: Body
	*/
	Payload *models.WorkloadCluster `json:"body,omitempty"`
}

// NewCreateWorkloadClusterOK creates CreateWorkloadClusterOK with default headers values
func NewCreateWorkloadClusterOK() *CreateWorkloadClusterOK {

	return &CreateWorkloadClusterOK{}
}

// WithPayload adds the payload to the create workload cluster o k response
func (o *CreateWorkloadClusterOK) WithPayload(payload *models.WorkloadCluster) *CreateWorkloadClusterOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create workload cluster o k response
func (o *CreateWorkloadClusterOK) SetPayload(payload *models.WorkloadCluster) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateWorkloadClusterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateWorkloadClusterBadRequestCode is the HTTP code returned for type CreateWorkloadClusterBadRequest
const CreateWorkloadClusterBadRequestCode int = 400

/*CreateWorkloadClusterBadRequest Bad request

swagger:response createWorkloadClusterBadRequest
*/
type CreateWorkloadClusterBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateWorkloadClusterBadRequest creates CreateWorkloadClusterBadRequest with default headers values
func NewCreateWorkloadClusterBadRequest() *CreateWorkloadClusterBadRequest {

	return &CreateWorkloadClusterBadRequest{}
}

// WithPayload adds the payload to the create workload cluster bad request response
func (o *CreateWorkloadClusterBadRequest) WithPayload(payload *models.Error) *CreateWorkloadClusterBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create workload cluster bad request response
func (o *CreateWorkloadClusterBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateWorkloadClusterBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateWorkloadClusterUnauthorizedCode is the HTTP code returned for type CreateWorkloadClusterUnauthorized
const CreateWorkloadClusterUnauthorizedCode int = 401

/*CreateWorkloadClusterUnauthorized Incorrect credentials

swagger:response createWorkloadClusterUnauthorized
*/
type CreateWorkloadClusterUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateWorkloadClusterUnauthorized creates CreateWorkloadClusterUnauthorized with default headers values
func NewCreateWorkloadClusterUnauthorized() *CreateWorkloadClusterUnauthorized {

	return &CreateWorkloadClusterUnauthorized{}
}

// WithPayload adds the payload to the create workload cluster unauthorized response
func (o *CreateWorkloadClusterUnauthorized) WithPayload(payload *models.Error) *CreateWorkloadClusterUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create workload cluster unauthorized response
func (o *CreateWorkloadClusterUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateWorkloadClusterUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateWorkloadClusterInternalServerErrorCode is the HTTP code returned for type CreateWorkloadClusterInternalServerError
const CreateWorkloadClusterInternalServerErrorCode int = 500

/*CreateWorkloadClusterInternalServerError Internal server error

swagger:response createWorkloadClusterInternalServerError
*/
type CreateWorkloadClusterInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateWorkloadClusterInternalServerError creates CreateWorkloadClusterInternalServerError with default headers values
func NewCreateWorkloadClusterInternalServerError() *CreateWorkloadClusterInternalServerError {

	return &CreateWorkloadClusterInternalServerError{}
}

// WithPayload adds the payload to the create workload cluster internal server error response
func (o *CreateWorkloadClusterInternalServerError) WithPayload(payload *models.Error) *CreateWorkloadClusterInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create workload cluster internal server error response
func (o *CreateWorkloadClusterInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateWorkloadClusterInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
