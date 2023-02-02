// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"github.com/ghodss/yaml"
)

const (
	// MinimumNumberOfParameters is 3
	MinimumNumberOfParameters int = 4
)

func main() {
	if len(os.Args) < MinimumNumberOfParameters {
		fmt.Println("error: missing argument")
		fmt.Printf("usage: %s <generator-name> <template-file> <values-file>\n", os.Args[0])
		os.Exit(1)
	}

	generatorName := os.Args[1]
	tmplPath := os.Args[2]
	valuesPath := os.Args[3]
	version := os.Args[4]

	tmpl := template.Must(template.ParseFiles(tmplPath))
	tmplData := buildTmplData(generatorName, valuesPath, version)
	err := tmpl.Execute(os.Stdout, tmplData)
	checkErr(err, "execute template")
}

type Data struct {
	GeneratorName       string
	Date                string
	Values              []*Value
	BundleDirForVersion string
}

type Value struct {
	ID          string
	Required    bool
	Description string
	Default     string
}

func buildTmplData(generatorName, valuesPath, version string) *Data {
	return &Data{
		GeneratorName:       generatorName,
		Date:                time.Now().Format(time.RFC1123),
		Values:              getValues(valuesPath),
		BundleDirForVersion: filepath.Join("addons", "packages", "pinniped", version, "bundle"),
	}
}

func getValues(valuesPath string) []*Value {
	valuesData, err := os.ReadFile(valuesPath)
	checkErr(err, "read values file")

	var valuesYAML map[string]interface{}
	err = yaml.Unmarshal(valuesData, &valuesYAML)
	checkErr(err, "unmarshal values YAML")

	values := buildValues(valuesYAML, "")
	sort.Slice(values, func(i, j int) bool {
		return values[i].ID < values[j].ID
	})
	return values
}

func buildValues(valuesYAML map[string]interface{}, baseKey string) []*Value {
	var values []*Value
	for key, value := range valuesYAML {
		id := baseKey + key
		if _, ok := value.(map[string]interface{}); ok {
			values = append(values, buildValues(value.(map[string]interface{}), id+".")...)
		} else {
			values = append(values, &Value{ID: id, Default: fmt.Sprintf("%v", value)})
		}
	}
	return values
}

func checkErr(err error, reason string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s: %s\n", reason, err.Error())
		os.Exit(1)
	}
}
