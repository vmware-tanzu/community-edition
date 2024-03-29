// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"

	types "github.com/docker/docker/api/types"
)

// DockerClient is an autogenerated mock type for the DockerClient type
type DockerClient struct {
	mock.Mock
}

// ContainerList provides a mock function with given fields: _a0, _a1
func (_m *DockerClient) ContainerList(_a0 context.Context, _a1 types.ContainerListOptions) ([]types.Container, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []types.Container
	if rf, ok := ret.Get(0).(func(context.Context, types.ContainerListOptions) []types.Container); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Container)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.ContainerListOptions) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerRemove provides a mock function with given fields: _a0, _a1, _a2
func (_m *DockerClient) ContainerRemove(_a0 context.Context, _a1 string, _a2 types.ContainerRemoveOptions) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ContainerRemoveOptions) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerStop provides a mock function with given fields: _a0, _a1, _a2
func (_m *DockerClient) ContainerStop(_a0 context.Context, _a1 string, _a2 *time.Duration) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *time.Duration) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Info provides a mock function with given fields: _a0
func (_m *DockerClient) Info(_a0 context.Context) (types.Info, error) {
	ret := _m.Called(_a0)

	var r0 types.Info
	if rf, ok := ret.Get(0).(func(context.Context) types.Info); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.Info)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewDockerInterfaceT interface {
	mock.TestingT
	Cleanup(func())
}

// NewDockerClient creates a new instance of DockerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDockerClient(t NewDockerInterfaceT) *DockerClient {
	mock := &DockerClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
