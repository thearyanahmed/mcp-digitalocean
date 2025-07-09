package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"github.com/invopop/jsonschema"
)

//go:generate go run .

// We generate the JSON schema from godo structs for CreateClusterRequest and CreateNodePoolRequest.
// This is necessary since we need to pass the NodePoolSpecs to the mcp tool as a raw argument.
// Ideally, we shouldn't have to copy the godo files around. However, it's currently not possible to without preserving the struct comments.
func main() {
	reflect := jsonschema.Reflector{
		BaseSchemaID:               "",
		Anonymous:                  true,
		AssignAnchor:               false,
		AllowAdditionalProperties:  true,
		RequiredFromJSONSchemaTags: true,
		DoNotReference:             true,
		ExpandedStruct:             true,
		FieldNameTag:               "",
	}

	err := reflect.AddGoComments("github.com/digitalocean/godo", "./")
	if err != nil {
		panic(fmt.Errorf("failed to add Go comments: %w", err))
	}

	createSchema, err := reflect.Reflect(&godo.KubernetesClusterCreateRequest{}).MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("failed to marshal cluster create schema: %w", err))
	}

	var createSchemaJSON bytes.Buffer
	if err := json.Indent(&createSchemaJSON, createSchema, "", "  "); err != nil {
		panic(fmt.Errorf("failed to indent JSON: %w", err))
	}

	// now write the schema to a file
	err = os.WriteFile("./cluster-create-schema.json", createSchemaJSON.Bytes(), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write schema to file: %w", err))
	}

	fmt.Println("Schema successfully written to cluster-create-schema.json")

	// Generate schema for KubernetesNodePoolCreateRequest
	npCreateSchema, err := reflect.Reflect(&godo.KubernetesNodePoolCreateRequest{}).MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("failed to marshal node pool create schema: %w", err))
	}

	// Prettify the JSON
	var npCreateSchemaJSON bytes.Buffer
	if err := json.Indent(&npCreateSchemaJSON, npCreateSchema, "", "  "); err != nil {
		panic(fmt.Errorf("failed to indent JSON: %w", err))
	}

	// Write the schema to a file
	err = os.WriteFile("./node-pool-create-schema.json", npCreateSchemaJSON.Bytes(), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write schema to file: %w", err))
	}

	fmt.Println("Update schema successfully written to node-pool-create-schema.json")
}
