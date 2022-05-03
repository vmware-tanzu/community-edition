// Code generated by go-swagger; DO NOT EDIT.

package cluster

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
)

// GetWorkloadClustersOKCode is the HTTP code returned for type GetWorkloadClustersOK
const GetWorkloadClustersOKCode int = 200

/*GetWorkloadClustersOK List of workload clusters being managed.

swagger:response getWorkloadClustersOK
*/
type GetWorkloadClustersOK struct {

	/*A list of workload clusters.
	  In: Body
	*/
	Payload []*models.WorkloadCluster `json:"body,omitempty"`
}

// NewGetWorkloadClustersOK creates GetWorkloadClustersOK with default headers values
func NewGetWorkloadClustersOK() *GetWorkloadClustersOK {

	return &GetWorkloadClustersOK{}
}

// WithPayload adds the payload to the get workload clusters o k response
func (o *GetWorkloadClustersOK) WithPayload(payload []*models.WorkloadCluster) *GetWorkloadClustersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get workload clusters o k response
func (o *GetWorkloadClustersOK) SetPayload(payload []*models.WorkloadCluster) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetWorkloadClustersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.WorkloadCluster, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetWorkloadClustersBadRequestCode is the HTTP code returned for type GetWorkloadClustersBadRequest
const GetWorkloadClustersBadRequestCode int = 400

/*GetWorkloadClustersBadRequest Bad request

swagger:response getWorkloadClustersBadRequest
*/
type GetWorkloadClustersBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetWorkloadClustersBadRequest creates GetWorkloadClustersBadRequest with default headers values
func NewGetWorkloadClustersBadRequest() *GetWorkloadClustersBadRequest {

	return &GetWorkloadClustersBadRequest{}
}

// WithPayload adds the payload to the get workload clusters bad request response
func (o *GetWorkloadClustersBadRequest) WithPayload(payload *models.Error) *GetWorkloadClustersBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get workload clusters bad request response
func (o *GetWorkloadClustersBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetWorkloadClustersBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetWorkloadClustersUnauthorizedCode is the HTTP code returned for type GetWorkloadClustersUnauthorized
const GetWorkloadClustersUnauthorizedCode int = 401

/*GetWorkloadClustersUnauthorized Incorrect credentials

swagger:response getWorkloadClustersUnauthorized
*/
type GetWorkloadClustersUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetWorkloadClustersUnauthorized creates GetWorkloadClustersUnauthorized with default headers values
func NewGetWorkloadClustersUnauthorized() *GetWorkloadClustersUnauthorized {

	return &GetWorkloadClustersUnauthorized{}
}

// WithPayload adds the payload to the get workload clusters unauthorized response
func (o *GetWorkloadClustersUnauthorized) WithPayload(payload *models.Error) *GetWorkloadClustersUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get workload clusters unauthorized response
func (o *GetWorkloadClustersUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetWorkloadClustersUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetWorkloadClustersInternalServerErrorCode is the HTTP code returned for type GetWorkloadClustersInternalServerError
const GetWorkloadClustersInternalServerErrorCode int = 500

/*GetWorkloadClustersInternalServerError Internal server error

swagger:response getWorkloadClustersInternalServerError
*/
type GetWorkloadClustersInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetWorkloadClustersInternalServerError creates GetWorkloadClustersInternalServerError with default headers values
func NewGetWorkloadClustersInternalServerError() *GetWorkloadClustersInternalServerError {

	return &GetWorkloadClustersInternalServerError{}
}

// WithPayload adds the payload to the get workload clusters internal server error response
func (o *GetWorkloadClustersInternalServerError) WithPayload(payload *models.Error) *GetWorkloadClustersInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get workload clusters internal server error response
func (o *GetWorkloadClustersInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetWorkloadClustersInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
