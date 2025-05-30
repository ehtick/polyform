package repeat

import (
	"github.com/EliCDavis/polyform/generator"
	"github.com/EliCDavis/polyform/nodes"
	"github.com/EliCDavis/polyform/refutil"
)

func init() {
	factory := &refutil.TypeFactory{}
	refutil.RegisterType[MeshNode](factory)
	refutil.RegisterType[CircleNode](factory)
	refutil.RegisterType[SplineNode](factory)
	refutil.RegisterType[LineNode](factory)
	refutil.RegisterType[FibonacciSphereNode](factory)
	refutil.RegisterType[FibonacciSpiralNode](factory)
	refutil.RegisterType[TRSNode](factory)
	refutil.RegisterType[TransformationNode](factory)
	refutil.RegisterType[nodes.Struct[SampleMeshSurfaceNode]](factory)
	refutil.RegisterType[nodes.Struct[polygonNode]](factory)
	refutil.RegisterType[nodes.Struct[GridNode]](factory)
	generator.RegisterTypes(factory)
}
