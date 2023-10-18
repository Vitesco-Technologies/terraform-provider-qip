/*
Copyright 2023 Vitesco Technologies Group AG

SPDX-License-Identifier: Apache-2.0

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/orderedmap"
)

var (
	typeNames   = flag.String("type", "", "comma-separated list of type names. required")
	packageName = flag.String("package", "", "name of the golang package. required")
	output      = flag.String("output", "", "output file name; default srcdir/<type>_type.go")
)

var uppercaseWords = []string{
	"ttl",
	"id",
	"rr",
}

const fileReadWrite = 0o644

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of qip_type:\n")
	fmt.Fprintf(os.Stderr, "  qip_type [flags] -type T\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("qip_type: ")

	flag.Usage = Usage
	flag.Parse()

	if len(*typeNames) == 0 || len(*packageName) == 0 {
		flag.Usage()
		os.Exit(2) //nolint:gomnd
	}

	types := strings.Split(*typeNames, ",")

	var g Generator

	fmt.Fprintf(&g.buf, "// Code generated by \"qip_type %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	fmt.Fprintf(&g.buf, "\n")
	fmt.Fprintf(&g.buf, "package %s", *packageName)
	fmt.Fprintf(&g.buf, "\n")

	for _, typeName := range types {
		g.Parse(typeName)
	}

	src := g.Format()

	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_type.go", types[0])
		// Output will be to local directory
		dir, _ := os.Getwd()
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}

	err := os.WriteFile(outputName, src, fileReadWrite)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Parse(typeName string) {
	name := strings.ToLower(typeName)
	schemaFile := name + ".json"

	if _, err := os.Stat(schemaFile); err != nil {
		log.Printf("warning: could not find schema file %s: %s", schemaFile, err)

		return
	}

	content, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Printf("warning: could not read schema file %s: %s", schemaFile, err)

		return
	}

	schema := orderedmap.New()

	err = json.Unmarshal(content, &schema)
	if err != nil {
		log.Printf("warning: could not parse JSON of schema file %s: %s", schemaFile, err)

		return
	}

	g.buf.Write(parseOrderedMap(typeName, schema).Bytes())
}

func parseOrderedMap(name string, data *orderedmap.OrderedMap) *bytes.Buffer {
	var buf bytes.Buffer

	subTypes := orderedmap.New()

	fmt.Fprintf(&buf, "type %s struct {\n", name)

	for _, key := range data.Keys() {
		field := key

		// Uppercase known prefixes or values.
		for _, word := range uppercaseWords {
			if strings.HasPrefix(field, word) {
				l := len(word)
				field = strings.ToUpper(field[0:l]) + field[l:]
			}
		}

		// Ensure field is public, by upper casing first character
		field = strings.ToUpper(field[0:1]) + field[1:]

		// parse types from the schema
		var goType string

		value, _ := data.Get(key)
		switch t := value.(type) {
		case string:
			// Most strings just contain the word "string" or exemplary values.
			// Some strings can refer to ENUMs, for now we are treating them all as strings.
			goType = "string"
		case float64:
			// This can be a float64 or int, QIP API refers to this as "integer" and "number".
			// As schema is unclear on this, we will apply it based on the field name.
			if strings.Contains(key, "Percent") {
				goType = "float64"
			} else {
				goType = "int"
			}
		case bool:
			goType = "bool"
		case orderedmap.OrderedMap:
			// Remember map to post process
			goType = name + field
			subTypes.Set(goType, value)
		case []interface{}:
			// An array with a single entry is another schema example
			if len(t) == 1 {
				if v, _ := t[0].(string); v == "string" {
					// Check if the array has a single element "string", then simple list of names
					goType = "[]string"
				} else if v, ok := t[0].(orderedmap.OrderedMap); ok {
					// Is a list of sub types
					goType = name + field
					subTypes.Set(goType, v)
					goType = "[]" + goType
				}
			}
		}

		if goType == "" {
			log.Printf("warning: unsupported type: %s is %T", key, value)

			if s, ok := value.(string); ok {
				log.Printf("warning: unsupported string type: %s is %s", key, s)
			}
		} else {
			fmt.Fprintf(&buf, "\t%s %s `json:\"%s,omitempty\"`\n", field, goType, key)
		}
	}

	buf.WriteString("}\n\n")

	// Process sub types
	for _, key := range subTypes.Keys() {
		if interf, ok := subTypes.Get(key); ok {
			if value, ok := interf.(orderedmap.OrderedMap); ok {
				buf.Write(parseOrderedMap(key, &value).Bytes())
			}
		}
	}

	return &buf
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) Format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")

		return g.buf.Bytes()
	}

	return src
}
