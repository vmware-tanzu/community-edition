// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/helloeave/json"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

const colWidth = 300
const indentation = `  `

type OutputWriter interface {
	SetKeys(headerKeys ...string)
	AddRow(items ...interface{})
	Render()
}

type OutputType string

const (
	TableOutputType = "table"
	YAMLOutputType  = "yaml"
	JSONOutputType  = "json"
)

// NewOutputWriter gets a new instance of our table output writer.
func NewOutputWriter(output io.Writer, outputFormat string, headers ...string) OutputWriter {
	// Initialize the output writer that we use under the covers
	ow := &outputwriter{}
	ow.out = output
	ow.outputFormat = outputFormat
	ow.keys = headers

	return ow
}

// outputwriter is our internal implementation.
type outputwriter struct {
	out          io.Writer
	keys         []string
	values       [][]string
	outputFormat string
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
	switch strings.ToLower(ow.outputFormat) {
	case JSONOutputType:
		renderJSON(ow)
	case YAMLOutputType:
		renderYAML(ow)
	default:
		renderTable(ow)
	}

	// ensures a break line after
	fmt.Fprintln(ow.out)
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
			item[keys[i]] = value
		}
		data = append(data, item)
	}

	return data
}

// renderJSON prints output as json
func renderJSON(ow *outputwriter) {
	data := ow.dataStruct()
	bytesJSON, err := json.MarshalSafeCollections(data)
	if err != nil {
		fmt.Fprint(ow.out, err)
		return
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, bytesJSON, "", indentation)
	if err != nil {
		fmt.Fprint(ow.out, err)
		return
	}

	fmt.Fprintf(ow.out, "%v", prettyJSON.String())
}

// renderYAML prints output as yaml
func renderYAML(ow *outputwriter) {
	data := ow.dataStruct()
	yamlInBytes, err := yaml.Marshal(data)
	if err != nil {
		fmt.Fprint(ow.out, err)
		return
	}

	fmt.Fprintf(ow.out, "%s", yamlInBytes)
}

// renderTable prints output as a table
func renderTable(ow *outputwriter) {
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
