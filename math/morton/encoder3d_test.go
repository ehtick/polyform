package morton_test

import (
	"testing"

	"github.com/EliCDavis/polyform/math/geometry"
	"github.com/EliCDavis/polyform/math/morton"
	"github.com/EliCDavis/vector/vector3"
	"github.com/stretchr/testify/assert"
)

func TestEncoder3D(t *testing.T) {

	tests := map[string]struct {
		Position vector3.Float64
		Encoder  morton.Encoder3D
		Delta    float64
		Encoded  uint64
	}{
		"Bottom Back Left": {
			Position: vector3.New(-1, -1, -1.),
			Encoder: morton.Encoder3D{
				Bounds:     geometry.NewAABB(vector3.Float64{}, vector3.Fill(2.)),
				Resolution: 2,
			},
			Delta:   0,
			Encoded: 0,
		},
		"Top Right Forward": {
			Position: vector3.New(1, 1, 1.),
			Encoder: morton.Encoder3D{
				Bounds:     geometry.NewAABB(vector3.Float64{}, vector3.Fill(2.)),
				Resolution: 2,
			},
			Delta:   0,
			Encoded: 0b111111,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			encoded := tc.Encoder.Encode(tc.Position)
			back := tc.Encoder.Decode(encoded)

			assert.Equal(t, tc.Encoded, encoded, "Encoded Index")

			assert.InDelta(t, tc.Position.X(), back.X(), tc.Delta, "Decoded X Axis")
			assert.InDelta(t, tc.Position.Y(), back.Y(), tc.Delta, "Decoded Y Axis")
			assert.InDelta(t, tc.Position.Z(), back.Z(), tc.Delta, "Decoded Z Axis")
		})
	}

}

func TestEncoder3D_Array(t *testing.T) {

	tests := map[string]struct {
		Positions []vector3.Float64
		Encoder   morton.Encoder3D
		Delta     float64
		Encoded   []uint64
	}{
		"Bottom Back Left": {
			Positions: []vector3.Float64{vector3.New(-1, -1, -1.)},
			Encoder: morton.Encoder3D{
				Bounds:     geometry.NewAABB(vector3.Float64{}, vector3.Fill(2.)),
				Resolution: 2,
			},
			Delta:   0,
			Encoded: []uint64{0},
		},
		"Top Right Forward": {
			Positions: []vector3.Float64{vector3.New(1, 1, 1.)},
			Encoder: morton.Encoder3D{
				Bounds:     geometry.NewAABB(vector3.Float64{}, vector3.Fill(2.)),
				Resolution: 2,
			},
			Delta:   0,
			Encoded: []uint64{0b111111},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			encoded := tc.Encoder.EncodeArray(tc.Positions)
			back := tc.Encoder.DecodeArray(encoded)
			for i, p := range tc.Positions {
				assert.Equal(t, tc.Encoded[i], encoded[i], "Encoded Index %d", i)

				assert.InDelta(t, p.X(), back[i].X(), tc.Delta, "Decoded X Axis %d", i)
				assert.InDelta(t, p.Y(), back[i].Y(), tc.Delta, "Decoded Y Axis %d", i)
				assert.InDelta(t, p.Z(), back[i].Z(), tc.Delta, "Decoded Z Axis %d", i)
			}

		})
	}

}

func TestExpand21(t *testing.T) {

	tests := map[string]struct {
		In       uint64
		Expected uint64
	}{
		"0b0": {},
		"0b1": {
			In:       uint64(0b1),
			Expected: uint64(0b1),
		},
		"0b11": {
			In:       uint64(0b11),
			Expected: uint64(0b1001),
		},
		"0b111": {
			In:       uint64(0b111),
			Expected: uint64(0b1001001),
		},
		"0b101": {
			In:       uint64(0b101),
			Expected: uint64(0b1000001),
		},
		"0b1111_1111": {
			In:       uint64(0b1111_1111),
			Expected: uint64(0b001001001001001001001001),
		},
		"0b1111_1111_1111": {
			In:       uint64(0b1111_1111_111),
			Expected: uint64(0b001001001001001001001001001001001),
		},
		"0b1111_1111_1111_1111_1111": {
			In:       uint64(0b1111_1111_1111_1111_1111),
			Expected: uint64(0b001001001001001001001001001001001001001001001001001001001001),
		},
		"0b1111_1111_1111_1111_1111_1": {
			In:       uint64(0b1111_1111_1111_1111_1111_1),
			Expected: uint64(0b001001001001001001001001001001001001001001001001001001001001001),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			encoded := morton.ExpandBits21(tc.In)
			back := morton.CompactBits21(encoded)

			assert.Equal(t, tc.Expected, encoded, "Encoded Index")
			assert.Equal(t, tc.In, back, "Decoded Index")
		})
	}

}
