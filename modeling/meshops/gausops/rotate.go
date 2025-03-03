package gausops

import (
	"github.com/EliCDavis/polyform/math/quaternion"
	"github.com/EliCDavis/polyform/modeling"
	"github.com/EliCDavis/polyform/nodes"
	"github.com/EliCDavis/vector/vector3"
	"github.com/EliCDavis/vector/vector4"
)

// https://github.com/aras-p/UnityGaussianSplatting/blob/ff268cfc6e12b4db80e2b1e9f14f7e31a68a8e25/package/Shaders/SplatUtilities.compute#L548
func RotateAttribute(m modeling.Mesh, attribute string, amount quaternion.Quaternion) modeling.Mesh {
	// q := quaternion.FromTheta(math.Pi, vector3.Forward[float64]())
	oldData := m.Float4Attribute(attribute)
	rotatedData := make([]vector4.Float64, oldData.Len())
	for i := 0; i < oldData.Len(); i++ {
		old := oldData.At(i)

		rot := amount.Normalize().Multiply(quaternion.New(vector3.New(old.Y(), old.Z(), old.W()), old.X()))
		rotatedData[i] = vector4.New(rot.W(), rot.Dir().X(), rot.Dir().Y(), rot.Dir().Z())

		// rot = amount.Multiply(quaternion.New(vector3.New(old.X(), old.Y(), old.Z()), old.W())).Normalize()
		// rotatedData[i] = vector4.New(rot.Dir().X(), rot.Dir().Y(), rot.Dir().Z(), rot.W())
	}

	return m.SetFloat4Attribute(attribute, rotatedData)
}

type RotateAttributeNode = nodes.Struct[modeling.Mesh, RotateAttributeNodeData]

type RotateAttributeNodeData struct {
	Mesh      nodes.NodeOutput[modeling.Mesh]
	Attribute nodes.NodeOutput[string]
	Amount    nodes.NodeOutput[quaternion.Quaternion]
}

func (rand RotateAttributeNodeData) Process() (modeling.Mesh, error) {
	if rand.Mesh == nil {
		return modeling.EmptyPointcloud(), nil
	}

	if rand.Amount == nil {
		return rand.Mesh.Value(), nil
	}

	attr := modeling.RotationAttribute
	if rand.Attribute != nil {
		attr = rand.Attribute.Value()
	}

	return RotateAttribute(rand.Mesh.Value(), attr, rand.Amount.Value()), nil
}
