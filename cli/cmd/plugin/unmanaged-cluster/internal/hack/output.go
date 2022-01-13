// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// The hack package contains temporary workarounds and other things that should
// be eventually removed.
package hack

// !!!!!!!!!!!!!!
// TODO(smcginnis) remove this once we can pick up the tanzu-framework version
// that includes the removal of bold headers in the table output.
// https://github.com/vmware-tanzu/tanzu-framework/commit/ddfbeceb775fb6bea8b6305787162b4614203bd0
// !!!!!!!!!!!!!!

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

const colWidth = 300
const indentation = `  `

// OutputWriter is an interface for something that can write output.
type OutputWriter interface {
	SetKeys(headerKeys ...string)
	AddRow(items ...interface{})
	Render()
}

// OutputType defines the format of the output desired.
type OutputType string

const (
	// TableOutputType specifies output should be in table format.
	TableOutputType OutputType = "table"
	// YAMLOutputType specifies output should be in yaml format.
	YAMLOutputType OutputType = "yaml"
	// JSONOutputType specifies output should be in json format.
	JSONOutputType OutputType = "json"
	// ListTableOutputType specified output should be in a list table format.
	ListTableOutputType OutputType = "listtable"
)

// outputwriter is our internal implementation.
type outputwriter struct {
	out          io.Writer
	keys         []string
	values       [][]string
	outputFormat OutputType
}

// NewOutputWriter gets a new instance of our output writer.
func NewOutputWriter(output io.Writer, outputFormat string, headers ...string) OutputWriter {
	// Initialize the output writer that we use under the covers
	ow := &outputwriter{}
	ow.out = output
	ow.outputFormat = OutputType(outputFormat)
	ow.keys = headers

	return ow
}

// SetKeys sets the values to use as the keys for the output values.
func (ow *outputwriter) SetKeys(headerKeys ...string) {
	// Overwrite whatever was used in initialization
	ow.keys = headerKeys
}

// AddRow appends a new row to our table.
func (ow *outputwriter) AddRow(items ...interface{}) {
	row := []string{}

	// Make sure all values are ultimately strings
	for _, item := range items {
		row = append(row, fmt.Sprintf("%v", item))
	}
	ow.values = append(ow.values, row)
}

// Render emits the generated table to the output once ready
func (ow *outputwriter) Render() {
	switch ow.outputFormat {
	case JSONOutputType:
		renderJSON(ow.out, ow.dataStruct())
	case YAMLOutputType:
		renderYAML(ow.out, ow.dataStruct())
	case ListTableOutputType:
		renderListTable(ow)
	default:
		renderTable(ow)
	}
}

func (ow *outputwriter) dataStruct() []map[string]string {
	data := []map[string]string{}
	keys := ow.keys
	for i, k := range keys {
		keys[i] = strings.ToLower(strings.ReplaceAll(k, " ", "_"))
	}

	for _, itemValues := range ow.values {
		item := map[string]string{}
		for i, value := range itemValues {
			if i == len(keys) {
				continue
			}
			item[keys[i]] = value
		}
		data = append(data, item)
	}

	return data
}

// objectwriter is our internal implementation.
type objectwriter struct {
	out          io.Writer
	data         interface{}
	outputFormat OutputType
}

// NewObjectWriter gets a new instance of our output writer.
func NewObjectWriter(output io.Writer, outputFormat string, data interface{}) OutputWriter {
	// Initialize the output writer that we use under the covers
	obw := &objectwriter{}
	obw.out = output
	obw.data = data
	obw.outputFormat = OutputType(outputFormat)

	return obw
}

// SetKeys sets the values to use as the keys for the output values.
func (obw *objectwriter) SetKeys(headerKeys ...string) {
	// Object writer does not have the concept of keys
	fmt.Fprintln(obw.out, "Programming error, attempt to add headers to object output")
}

// AddRow appends a new row to our table.
func (obw *objectwriter) AddRow(items ...interface{}) {
	// Object writer does not have the concept of keys
	fmt.Fprintln(obw.out, "Programming error, attempt to add rows to object output")
}

// Render emits the generated table to the output once ready
func (obw *objectwriter) Render() {
	switch obw.outputFormat {
	case JSONOutputType:
		renderJSON(obw.out, obw.data)
	case YAMLOutputType:
		renderYAML(obw.out, obw.data)
	default:
		fmt.Fprintf(obw.out, "Invalid output format: %v\n", obw.outputFormat)
	}
}

// renderJSON prints output as json
func renderJSON(out io.Writer, data interface{}) {
	bytesJSON, err := json.MarshalIndent(data, "", indentation)
	if err != nil {
		fmt.Fprint(out, err)
		return
	}

	fmt.Fprintf(out, "%v", string(bytesJSON))
	fmt.Println()
}

// renderYAML prints output as yaml
func renderYAML(out io.Writer, data interface{}) {
	yamlInBytes, err := yaml.Marshal(data)
	if err != nil {
		fmt.Fprint(out, err)
		return
	}

	fmt.Fprintf(out, "%s", yamlInBytes)
}

// renderListTable prints output as a list table.
func renderListTable(ow *outputwriter) {
	headerLength := 10
	for _, header := range ow.keys {
		length := len(header) + 2
		if length > headerLength {
			headerLength = length
		}
	}

	for i, header := range ow.keys {
		row := []string{}
		for _, data := range ow.values {
			if i >= len(data) {
				// There are more headers than values, leave it blank
				continue
			}
			row = append(row, data[i])
		}
		headerLabel := strings.ToUpper(header) + ":"
		values := strings.Join(row, ", ")
		fmt.Fprintf(ow.out, "%-"+strconv.Itoa(headerLength)+"s   %s\n", headerLabel, values)
	}
}

// renderTable prints output as a table
func renderTable(ow *outputwriter) {
	// Drop values if there aren't as many as the headers
	headerLength := len(ow.keys)
	for i, values := range ow.values {
		if len(values) <= headerLength {
			continue
		}

		ow.values[i] = values[:headerLength]
	}
	table := tablewriter.NewWriter(ow.out)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetColWidth(colWidth)
	table.SetTablePadding("\t\t")
	table.SetHeader(ow.keys)
	table.AppendBulk(ow.values)
	table.Render()
}
