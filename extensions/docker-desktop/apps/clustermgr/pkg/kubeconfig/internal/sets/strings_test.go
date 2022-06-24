// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package sets

import (
	"reflect"
	"testing"
)

func TestNewString(t *testing.T) {
	sVal := NewString("test", "test", "testing", "string1", "test", "string2", "string12", "string1")
	if sVal.Len() != 5 {
		t.Errorf("Expected the Length to be 5, got %d", sVal.Len())
	}
}

func TestStringKeySet(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		panic          bool
		expectedLength int
	}{
		{
			name: "Valid map",
			input: map[string]string{
				"k1": "v1",
				"k2": "v1",
				"k3": "v3",
				"k4": "v0",
			},
			panic:          false,
			expectedLength: 4,
		},
		{
			name:  "invalid non map type",
			input: []string{"test", "test"},
			panic: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected %s to panic, but it did not", tt.name)
					}
				}()
			}
			val := StringKeySet(tt.input)
			if !tt.panic {
				if val.Len() != tt.expectedLength {
					t.Errorf("len(StringKeySet(%v)) = %d, expected %d", tt.input, val.Len(), tt.expectedLength)
				}
			}
		})
	}
}

func TestStringUpdateOps(t *testing.T) {
	sVal := NewString("test", "test", "testing", "string1", "test", "string2", "string12", "string1")
	tests := []struct {
		name         string
		input        []string
		expected     []string
		functionName string
	}{
		{
			name:         "Insert With Duplicate",
			input:        []string{"abcd", "1234", "test"},
			expected:     []string{"test", "testing", "string1", "string2", "string12", "abcd", "1234"},
			functionName: "Insert",
		},
		{
			name:         "Delete Non existing with valid entries",
			input:        []string{"abcd", "1234", "test"},
			expected:     []string{"testing", "string1", "string2", "string12"},
			functionName: "Delete",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fHandler := reflect.ValueOf(sVal).MethodByName(tt.functionName)
			var input []reflect.Value
			for _, i := range tt.input {
				input = append(input, reflect.ValueOf(i))
			}
			returnVal := fHandler.Call(input)
			got := returnVal[0].Interface().(String)
			if !got.Equal(NewString(tt.expected...)) {
				t.Errorf("Expected %v, got %v", NewString(tt.expected...), got)
			}
		})
	}
}

func TestExistenceCheck(t *testing.T) {
	sVal := NewString("test", "test", "testing", "string1", "test", "string2", "string12", "string1")
	tests := []struct {
		name         string
		functionName string
		input        []string
		expected     bool
	}{
		{
			name:         "Should return true if exists",
			functionName: "Has",
			input:        []string{"test"},
			expected:     true,
		},
		{
			name:         "Should return false if not exists",
			functionName: "Has",
			input:        []string{"tests"},
			expected:     false,
		},
		{
			name:         "should return true if all items exist",
			functionName: "HasAll",
			input:        []string{"test", "string2"},
			expected:     true,
		},
		{
			name:         "should return false if not all items exist",
			functionName: "HasAll",
			input:        []string{"test", "string123"},
			expected:     false,
		},
		{
			name:         "should return true if any one item exists",
			functionName: "HasAny",
			input:        []string{"test", "string123"},
			expected:     true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fHandler := reflect.ValueOf(sVal).MethodByName(tt.functionName)
			var input []reflect.Value
			for _, i := range tt.input {
				input = append(input, reflect.ValueOf(i))
			}
			returnVal := fHandler.Call(input)
			got := returnVal[0].Interface().(bool)
			if got != tt.expected {
				t.Errorf("Expected %v, Got %v", tt.expected, got)
			}
		})
	}
}

func TestSetOperations(t *testing.T) {
	sVal := NewString("a1", "a2", "a3", "a5")
	tests := []struct {
		name         string
		input        []string
		expected     []string
		functionName string
	}{
		{
			name:         "Set Difference",
			input:        []string{"a1", "a2", "a3", "a4"},
			expected:     []string{"a5"},
			functionName: "Difference",
		},
		{
			name:         "Set Union",
			input:        []string{"a1", "a2", "a3", "a4"},
			expected:     []string{"a1", "a2", "a3", "a4", "a5"},
			functionName: "Union",
		},
		{
			name:         "Set Intersection",
			input:        []string{"a1", "a2", "a3", "a4"},
			expected:     []string{"a1", "a2", "a3"},
			functionName: "Intersection",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fHandler := reflect.ValueOf(sVal).MethodByName(tt.functionName)
			s2 := NewString(tt.input...)
			var input []reflect.Value
			input = append(input, reflect.ValueOf(s2))
			returnVal := fHandler.Call(input)
			got := returnVal[0].Interface().(String)
			if !got.Equal(NewString(tt.expected...)) {
				t.Errorf("Expected %v, got %v", NewString(tt.expected...), got)
			}
		})
	}
}

func TestSortableOperations(t *testing.T) {
	tests := []struct {
		name         string
		input        []string
		expected     []string
		functionName string
		count        int
	}{
		{
			name:         "Sorted String set",
			input:        []string{"1", "2", "10", "a"},
			expected:     []string{"1", "10", "2", "a"},
			functionName: "List",
		},
		{
			name:         "Sorted String set with special characters",
			input:        []string{"1", "2", "10", "a", "%"},
			expected:     []string{"%", "1", "10", "2", "a"},
			functionName: "List",
		},
		{
			name:         "Unsorted String set with special characters",
			input:        []string{"1", "2", "10", "a", "%", "1", "%"},
			expected:     []string{"1", "2", "10", "a", "%"},
			functionName: "UnsortedList",
			count:        5,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sVal := NewString(tt.input...)
			fHandler := reflect.ValueOf(sVal).MethodByName(tt.functionName)
			var input []reflect.Value
			returnVal := fHandler.Call(input)
			got := returnVal[0].Interface().([]string)
			if tt.count == 0 {
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, got)
				}
			} else {
				if len(got) != tt.count {
					t.Errorf("Expected %v, got %v", tt.expected, got)
				}
			}
		})
	}
}
