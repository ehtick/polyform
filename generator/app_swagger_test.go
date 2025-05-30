package generator_test

import (
	"testing"

	"github.com/EliCDavis/polyform/formats/ply"
	"github.com/EliCDavis/polyform/formats/swagger"
	"github.com/EliCDavis/polyform/generator"
	"github.com/EliCDavis/polyform/generator/manifest"
	"github.com/EliCDavis/polyform/generator/parameter"
	"github.com/EliCDavis/polyform/math/geometry"
	"github.com/EliCDavis/polyform/modeling"
	"github.com/EliCDavis/polyform/modeling/meshops"
	"github.com/EliCDavis/polyform/nodes"
	"github.com/EliCDavis/vector/vector3"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerFromGraph_SingleParameterSingleProducer(t *testing.T) {

	appName := "Test Graph"
	appVersion := "Test Graph Version"
	appDescription := "Test Graph Description"
	producerFileName := "test.txt"
	app := generator.App{
		Name:        appName,
		Version:     appVersion,
		Description: appDescription,
		Files: map[string]nodes.Output[manifest.Manifest]{
			producerFileName: buildTextArifact(&parameter.Value[string]{
				Name:         "Welp",
				DefaultValue: "yee",
			}),
		},
	}

	// ACT ====================================================================
	spec := app.SwaggerSpec()

	// ASSERT =================================================================
	assert.Equal(t, "2.0", spec.Version)

	// Info Details
	assert.Equal(t, appName, spec.Info.Title)
	assert.Equal(t, appDescription, spec.Info.Description)
	assert.Equal(t, appVersion, spec.Info.Version)

	// Producers
	assert.Len(t, spec.Paths, 1)
	assert.Contains(t, spec.Paths, "/producer/value/test.txt")

	path := spec.Paths["/producer/value/test.txt"]
	assert.Len(t, path, 1)
	assert.Contains(t, path, swagger.PostRequestMethod)

	request := path[swagger.PostRequestMethod]
	if assert.Len(t, request.Consumes, 1) {
		assert.Equal(t, "application/json", request.Consumes[0])
	}

	// Parameter
	parameters := request.Parameters
	assert.Len(t, parameters, 1)
	assert.Equal(t, swagger.BodyParameterLocation, parameters[0].In)
	assert.Equal(t, "Request", parameters[0].Name)
	assert.Equal(t, false, parameters[0].Required)
	assert.Equal(
		t,
		swagger.SchemaObject{
			Ref: "#/definitions/TestTxtRequest",
		},
		parameters[0].Schema,
	)

	// Response
	responses := request.Responses
	assert.Len(t, responses, 1)
	assert.Contains(t, responses, 200)

	// Success Response
	successResponse := responses[200]
	assert.Equal(t, "Producer Payload", successResponse.Description)
	// assert.Equal(
	// 	t,
	// 	swagger.SchemaObject{
	// 		Ref: "#/definitions/TestTxtRequest",
	// 	},
	// 	successResponse.Schema,
	// )

	// Definitions
	assert.Len(t, spec.Definitions, 1)
	assert.Contains(t, spec.Definitions, "TestTxtRequest")

	reqDef := spec.Definitions["TestTxtRequest"]
	assert.Equal(t, "object", reqDef.Type)
	assert.Len(t, reqDef.Required, 0)
	assert.Len(t, reqDef.Properties, 1)
	assert.Contains(t, reqDef.Properties, "Welp")

	prop := reqDef.Properties["Welp"]
	assert.Equal(t, swagger.StringPropertyType, prop.Type)
}

func TestSwaggerFromGraph_MultipleParametersSingleProducer(t *testing.T) {

	aabb := &parameter.AABB{
		Name:        "Bounding Box",
		Description: "Box to crop gaussian splat by",
	}

	vec := &parameter.Vector3{
		Name:        "Translation",
		Description: "Amount to shift the mesh by",
	}

	mesh := &meshops.TranslateAttribute3DNode{
		Data: meshops.TranslateAttribute3DNodeData{
			Amount: nodes.GetNodeOutputPort[vector3.Float64](vec, "Value"),
		},
	}

	crop := &meshops.CropAttribute3DNode{
		Data: meshops.CropAttribute3DNodeData{
			AABB: nodes.GetNodeOutputPort[geometry.AABB](aabb, "Value"),
			Mesh: nodes.GetNodeOutputPort[modeling.Mesh](mesh, "Out"),
		},
	}

	ply := &ply.ManifestNode{
		Data: ply.ManifestNodeData{
			Mesh: nodes.GetNodeOutputPort[modeling.Mesh](crop, "Out"),
		},
	}

	app := generator.App{
		Name:        "Example Graph",
		Version:     "1.0.0",
		Description: "Example graph that contains multiple parameters",
		Files: map[string]nodes.Output[manifest.Manifest]{
			"example.glb": nodes.GetNodeOutputPort[manifest.Manifest](ply, "Out"),
		},
	}

	// ACT ====================================================================
	spec := app.SwaggerSpec()

	// ASSERT =================================================================
	assert.Equal(t, "2.0", spec.Version)

	// Info Details
	assert.Equal(t, "Example Graph", spec.Info.Title)
	assert.Equal(t, "Example graph that contains multiple parameters", spec.Info.Description)
	assert.Equal(t, "1.0.0", spec.Info.Version)

	// Producers
	assert.Len(t, spec.Paths, 1)
	assert.Contains(t, spec.Paths, "/producer/value/example.glb")

	path := spec.Paths["/producer/value/example.glb"]
	assert.Len(t, path, 1)
	assert.Contains(t, path, swagger.PostRequestMethod)
	request := path[swagger.PostRequestMethod]

	// Parameter
	parameters := request.Parameters
	assert.Len(t, parameters, 1)
	assert.Equal(t, swagger.BodyParameterLocation, parameters[0].In)
	assert.Equal(t, "Request", parameters[0].Name)
	assert.Equal(t, false, parameters[0].Required)
	assert.Equal(
		t,
		swagger.SchemaObject{
			Ref: "#/definitions/ExampleGlbRequest",
		},
		parameters[0].Schema,
	)

	// Response
	responses := request.Responses
	assert.Len(t, responses, 1)
	assert.Contains(t, responses, 200)

	// Success Response
	successResponse := responses[200]
	assert.Equal(t, "Producer Payload", successResponse.Description)
	// assert.Equal(
	// 	t,
	// 	swagger.SchemaObject{
	// 		Ref: "#/definitions/TestTxtRequest",
	// 	},
	// 	successResponse.Schema,
	// )

	// Definitions
	assert.Len(t, spec.Definitions, 3)
	assert.Contains(t, spec.Definitions, "ExampleGlbRequest")

	reqDef := spec.Definitions["ExampleGlbRequest"]
	assert.Equal(t, "object", reqDef.Type)
	assert.Len(t, reqDef.Required, 0)
	assert.Len(t, reqDef.Properties, 2)
	assert.Contains(t, reqDef.Properties, "BoundingBox")
	assert.Contains(t, reqDef.Properties, "Translation")

	boxProp := reqDef.Properties["BoundingBox"]
	assert.Equal(t, swagger.PropertyType(""), boxProp.Type)
	assert.Equal(t, "#/definitions/AABB", boxProp.Ref)

	vectorProp := reqDef.Properties["Translation"]
	assert.Equal(t, swagger.PropertyType(""), vectorProp.Type)
	assert.Equal(t, "#/definitions/Vector3", vectorProp.Ref)
}
