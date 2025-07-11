package generator_test

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"

	"github.com/EliCDavis/polyform/generator"
	"github.com/EliCDavis/polyform/generator/manifest"
	"github.com/EliCDavis/polyform/generator/manifest/basics"
	"github.com/EliCDavis/polyform/generator/parameter"
	"github.com/EliCDavis/polyform/nodes"
	"github.com/stretchr/testify/assert"
)

func buildTextArifact(p *parameter.String) nodes.Output[manifest.Manifest] {
	return nodes.GetNodeOutputPort[manifest.Manifest](
		&basics.TextNode{
			Data: basics.TextNodeData{
				In: nodes.GetNodeOutputPort[string](p, "Value"),
			},
		},
		"Out",
	)
}

func TestGetAndApplyGraph(t *testing.T) {
	appName := "Test Graph"
	appVersion := "Test Graph"
	appDescription := "Test Graph"
	producerFileName := "test.txt"
	app := generator.App{
		Name:        appName,
		Version:     appVersion,
		Description: appDescription,
		Files: map[string]nodes.Output[manifest.Manifest]{
			producerFileName: buildTextArifact(&parameter.String{
				Name:         "Welp",
				DefaultValue: "yee",
			}),
		},
	}

	// ACT ====================================================================
	graphData := app.Schema()
	err := app.ApplySchema(graphData)
	graphAgain := app.Schema()

	// ASSERT =================================================================
	assert.NoError(t, err)
	assert.Equal(t, appName, app.Name)
	assert.Equal(t, appVersion, app.Version)
	assert.Equal(t, appDescription, app.Description)
	assert.Equal(t, string(graphData), string(graphAgain))
	b := &bytes.Buffer{}
	manifest := app.Files[producerFileName].Value()
	art := manifest.Entries[manifest.Main]
	err = art.Artifact.Write(b)
	assert.NoError(t, err)
	assert.Equal(t, "yee", b.String())
}

func TestAppCommand_Outline(t *testing.T) {
	appName := "Test Graph"
	appVersion := "Test Graph"
	appDescription := "Test Graph"
	producerFileName := "test.txt"

	outBuf := &bytes.Buffer{}

	app := generator.App{
		Name:        appName,
		Version:     appVersion,
		Description: appDescription,
		Files: map[string]nodes.Output[manifest.Manifest]{
			producerFileName: buildTextArifact(&parameter.String{
				Name:         "Welp",
				DefaultValue: "yee",
			}),
		},

		Out: outBuf,
	}

	// ACT ====================================================================
	err := app.Run([]string{"polyform", "outline"})
	contents, readErr := io.ReadAll(outBuf)

	// ASSERT =================================================================
	assert.NoError(t, err)
	assert.NoError(t, readErr)
	assert.Equal(t, `{
    "producers": {
        "test.txt": {
            "nodeID": "Node-1",
            "port": "Out"
        }
    },
    "nodes": {
        "Node-0": {
            "type": "github.com/EliCDavis/polyform/generator/parameter.Value[string]",
            "name": "Welp",
            "assignedInput": {},
            "output": {
                "Value": {
                    "version": 0
                }
            },
            "parameter": {
                "name": "Welp",
                "description": "",
                "type": "string",
                "defaultValue": "yee",
                "currentValue": "yee"
            }
        },
        "Node-1": {
            "type": "github.com/EliCDavis/polyform/nodes.Struct[github.com/EliCDavis/polyform/generator/manifest/basics.TextNodeData]",
            "name": "test.txt",
            "assignedInput": {
                "In": {
                    "id": "Node-0",
                    "port": "Value"
                }
            },
            "output": {
                "Out": {
                    "version": -1
                }
            }
        }
    },
    "notes": null,
    "variables": {
        "variables": {},
        "subgroups": {}
    }
}`, string(contents))
}

func TestAppCommand_Zip(t *testing.T) {
	appName := "Test Graph"
	appVersion := "Test Graph"
	appDescription := "Test Graph"
	producerFileName := "test"

	outBuf := &bytes.Buffer{}

	app := generator.App{
		Name:        appName,
		Version:     appVersion,
		Description: appDescription,
		Files: map[string]nodes.Output[manifest.Manifest]{
			producerFileName: buildTextArifact(&parameter.String{
				Name:         "Welp",
				DefaultValue: "yee",
			}),
		},
		Out: outBuf,
	}

	// ACT ====================================================================
	err := app.Run([]string{"polyform", "zip"})
	data := outBuf.Bytes()

	r, zipErr := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	// ASSERT =================================================================
	assert.NoError(t, err)
	assert.NoError(t, zipErr)
	assert.Len(t, r.File, 1)
	assert.Equal(t, "test/text.txt", r.File[0].Name)

	rc, err := r.File[0].Open()
	assert.NoError(t, err)

	buf, err := io.ReadAll(rc)
	assert.NoError(t, err)
	assert.Equal(t, "yee", string(buf))
}

func TestAppCommand_Swagger(t *testing.T) {
	appName := "Test Graph"
	appVersion := "Test Graph"
	appDescription := "Test Graph"
	producerFileName := "test.txt"

	outBuf := &bytes.Buffer{}

	app := generator.App{
		Name:        appName,
		Version:     appVersion,
		Description: appDescription,
		Files: map[string]nodes.Output[manifest.Manifest]{
			producerFileName: buildTextArifact(&parameter.String{
				Name:         "Welp",
				DefaultValue: "yee",
				Description:  "I'm a description",
			}),
		},

		Out: outBuf,
	}

	// ACT ====================================================================
	err := app.Run([]string{"polyform", "swagger"})
	contents, readErr := io.ReadAll(outBuf)

	// ASSERT =================================================================
	assert.NoError(t, err)
	assert.NoError(t, readErr)
	assert.Equal(t, `{
    "swagger": "2.0",
    "info": {
        "title": "Test Graph",
        "description": "Test Graph",
        "version": "Test Graph"
    },
    "paths": {
        "/producer/value/test.txt": {
            "post": {
                "summary": "",
                "description": "",
                "produces": [],
                "consumes": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "Producer Payload"
                    }
                },
                "parameters": [
                    {
                        "in": "body",
                        "name": "Request",
                        "schema": {
                            "$ref": "#/definitions/TestTxtRequest"
                        }
                    }
                ]
            }
        }
    },
    "definitions": {
        "TestTxtRequest": {
            "type": "object",
            "properties": {
                "Welp": {
                    "type": "string",
                    "description": "I'm a description"
                }
            }
        }
    }
}`, string(contents))
}

func TestAppCommand_New(t *testing.T) {
	appName := "Test Graph"
	appVersion := "Test Graph"
	appDescription := "Test Graph"
	producerFileName := "test.txt"

	outBuf := &bytes.Buffer{}

	app := generator.App{
		Name:        appName,
		Version:     appVersion,
		Description: appDescription,
		Files: map[string]nodes.Output[manifest.Manifest]{
			producerFileName: buildTextArifact(&parameter.String{
				Name:         "Welp",
				DefaultValue: "yee",
			}),
		},

		Out: outBuf,
	}

	// ACT ====================================================================
	err := app.Run([]string{
		"polyform", "new",
		"--name", "My New Graph",
		"--description", "This is just a test",
		"--version", "v1.0.2",
		"--author", "Test Runner",
	})
	contents, readErr := io.ReadAll(outBuf)

	// ASSERT =================================================================
	assert.NoError(t, err)
	assert.NoError(t, readErr)
	assert.Equal(t, `{
	"name": "My New Graph",
	"version": "v1.0.2",
	"description": "This is just a test",
	"authors": [
		{
			"name": "Test Runner"
		}
	],
	"producers": null,
	"nodes": null,
	"variables": {
		"variables": null,
		"subgroups": null
	}
}`, string(contents))
}
