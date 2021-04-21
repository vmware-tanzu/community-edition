// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const colWidth = 300

type TableOutputWriter interface {
	SetHeaders(headers ...string)
	AddRow(items ...interface{})
	Render()
}

// convertToUpper will make sure all entries are upper cased.
func convertToUpper(headers []string) []string {
	head := []string{}
	for _, item := range headers {
		head = append(head, strings.ToUpper(item))
	}

	return head
}

// NewTableWriter gets a new instance of our table output writer.
func NewTableWriter(output io.Writer, headers ...string) TableOutputWriter {
	// Initialize the output writer that we use under the covers
	table := tablewriter.NewWriter(output)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetColWidth(colWidth)
	table.SetTablePadding("\t\t")
	table.SetHeader(convertToUpper(headers))

	t := &tableoutputwriter{}
	t.out = output
	t.table = table

	return t
}

// tableoutputwriter is our internal implementation.
type tableoutputwriter struct {
	out   io.Writer
	table *tablewriter.Table
}

func (t *tableoutputwriter) SetHeaders(headers ...string) {
	// Overwrite whatever was used in initialization
	t.table.SetHeader(convertToUpper(headers))
}

// AddRow appends a new row to our table.
func (t *tableoutputwriter) AddRow(items ...interface{}) {
	row := []string{}

	// Make sure all values are ultimately strings
	for _, item := range items {
		row = append(row, fmt.Sprintf("%v", item))
	}
	t.table.Append(row)
}

// Render emits the generated table to the output once ready
func (t *tableoutputwriter) Render() {
	t.table.Render()

	// ensures a break line after we flush the tabwriter
	fmt.Fprintln(t.out)
}
