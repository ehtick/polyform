package nodes_test

// import (
// 	"testing"

// 	"github.com/EliCDavis/polyform/modeling"
// 	"github.com/EliCDavis/polyform/modeling/primitives"
// 	"github.com/EliCDavis/polyform/modeling/repeat"
// 	"github.com/EliCDavis/polyform/nodes"
// 	"github.com/stretchr/testify/assert"
// )

// type CombineNode = nodes.Struct[CombineData]

// type CombineData struct {
// 	Meshes []nodes.Output[modeling.Mesh]
// }

// func (cn CombineData) Out() nodes.StructOutput[modeling.Mesh] {
// 	finalMesh := modeling.EmptyMesh(modeling.TriangleTopology)

// 	for _, n := range cn.Meshes {
// 		finalMesh = finalMesh.Append(n.Value())
// 	}

// 	return finalMesh, nil
// }

// func TestNodes(t *testing.T) {

// 	times := nodes.NewValue(5)

// 	transforms := &repeat.CircleNode{
// 		Data: repeat.CircleNodeData{
// 			Radius: nodes.NewValue(15.),
// 			Times:  nodes.NewValue(5),
// 		},
// 	}

// 	repeated := &repeat.MeshNode{
// 		Data: repeat.MeshNodeData{
// 			Mesh: &repeat.MeshNode{
// 				Data: repeat.MeshNodeData{
// 					Mesh: nodes.NewValue(primitives.UVSphere(1, 10, 10)),
// 					Transforms: &repeat.CircleNode{
// 						Data: repeat.CircleNodeData{
// 							Radius: nodes.NewValue(5.),
// 							Times:  times,
// 						},
// 					},
// 				},
// 			},
// 			Transforms: transforms,
// 		},
// 	}

// 	transforms.
// 		Out().
// 		Node().
// 		SetInput("Times", nodes.TypedPort{Port: nodes.NewValue(30)})

// 	combined := CombineNode{
// 		Data: CombineData{
// 			Meshes: []nodes.Output[modeling.Mesh]{
// 				repeated.Out(),
// 				(&repeat.MeshNode{
// 					Data: repeat.MeshNodeData{
// 						Mesh: nodes.NewValue(primitives.UVSphere(1, 10, 10)),
// 						Transforms: &repeat.CircleNode{
// 							Data: repeat.CircleNodeData{
// 								Radius: nodes.NewValue(5.),
// 								Times:  times,
// 							},
// 						},
// 					},
// 				}).Out(),
// 			},
// 		},
// 	}

// 	combinedInputs := combined.Inputs()
// 	assert.Len(t, combinedInputs, 1)
// 	assert.Equal(t, "github.com/EliCDavis/polyform/modeling.Mesh", combinedInputs[0].Type, combinedInputs[0].Array)
// 	assert.True(t, combinedInputs[0].Array)

// 	combinedDeps := combined.Out().Node().Dependencies()
// 	assert.Len(t, combinedDeps, 2)

// 	// Stage changes
// 	out := combined.Out()

// 	out.Value()
// 	times.Set(13)
// 	out.Value()

// 	deps := repeated.Out().Node().Dependencies()
// 	assert.Len(t, deps, 2)
// 	// assert.Equal(t, []nodes.Output{{
// 	// 	// Name: "Out",
// 	// 	Type: "github.com/EliCDavis/polyform/modeling.Mesh",
// 	// }}, combined.Out().Node().Outputs())

// 	assert.Equal(t, "nodes_test", combined.Path())
// 	assert.Equal(t, "modeling/repeat", repeated.Path())
// 	// obj.Save("test.obj", repeat.Value())
// }
