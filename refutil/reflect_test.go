package refutil_test

import (
	"io"
	"testing"

	"github.com/EliCDavis/polyform/refutil"
	"github.com/EliCDavis/vector/vector3"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	A int
	B bool
	C vector3.Float64
	D *vector3.Float64
}

func (ts TestStruct) Read(b []byte) (int, error) {
	return 0, nil
}

func (ts TestStruct) ABC() error {
	return nil
}

func (ts TestStruct) XYZ() int {
	return 1
}

type genericTestStruct[T any] struct {
}

func (t genericTestStruct[T]) TypeWithPackage() string {
	var v T
	return refutil.GetTypeWithPackage(v)
}

func TestFuncNamesOfType(t *testing.T) {
	ts := TestStruct{}
	v := refutil.FuncNamesOfType[error](ts)

	assert.Len(t, v, 1)
	assert.Contains(t, v, "ABC")
}

func TestFuncNamesOfType_Interface(t *testing.T) {
	ts := TestStruct{}
	var reader io.Reader = &ts
	v := refutil.FuncNamesOfType[error](reader)

	assert.Len(t, v, 1)
	assert.Contains(t, v, "ABC")
}

func TestGenericFieldValuesOfType(t *testing.T) {
	ts := TestStruct{}

	v := refutil.GenericFieldTypes("vector3.Vector", ts)
	assert.Len(t, v, 1)
	assert.Equal(t, "float64", v["C"])

	v = refutil.GenericFieldTypes("vector3.Vector", &ts)
	assert.Len(t, v, 1)
	assert.Equal(t, "float64", v["C"])
}

func TestGetTypeWithPackageGeneric(t *testing.T) {
	assert.Equal(t, "int", genericTestStruct[int]{}.TypeWithPackage())
	assert.Equal(t, "github.com/EliCDavis/vector/vector3.Array[float64]", genericTestStruct[vector3.Array[float64]]{}.TypeWithPackage())
	assert.Equal(t, "github.com/EliCDavis/vector/vector3.Array[float64]", genericTestStruct[*vector3.Array[float64]]{}.TypeWithPackage())
	assert.Equal(t, "[]float64", genericTestStruct[[]float64]{}.TypeWithPackage())
}

func TestGetPackagePath(t *testing.T) {
	// var reader io.Reader
	tests := map[string]struct {
		input any
		want  string
	}{
		"nil": {
			input: nil,
			want:  "",
		},
		"string": {
			input: "test",
			want:  "",
		},
		"std lib": {
			input: io.Discard,
			want:  "io",
		},
		"external lib": {
			input: vector3.New(1, 2, 3),
			want:  "github.com/EliCDavis/vector/vector3",
		},
		"pointer external lib": {
			input: &vector3.Vector[float64]{},
			want:  "github.com/EliCDavis/vector/vector3",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := refutil.GetPackagePath(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetTypeWithPackage(t *testing.T) {
	// var reader io.Reader
	tests := map[string]struct {
		input any
		want  string
	}{
		"nil": {
			input: nil,
			want:  "nil",
		},
		"string": {
			input: "test",
			want:  "string",
		},
		"std lib": {
			input: io.Discard,
			want:  "io.discard",
		},
		"external lib": {
			input: vector3.New(1, 2, 3),
			want:  "github.com/EliCDavis/vector/vector3.Vector[int]",
		},
		"pointer external lib": {
			input: &vector3.Vector[float64]{},
			want:  "github.com/EliCDavis/vector/vector3.Vector[float64]",
		},
		"array external lib": {
			input: []vector3.Vector[float64]{},
			want:  "[]github.com/EliCDavis/vector/vector3.Vector[float64]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := refutil.GetTypeWithPackage(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetName(t *testing.T) {
	// var reader io.Reader
	var v *vector3.Vector[float64]
	tests := map[string]struct {
		input any
		want  string
	}{
		"nil": {
			input: nil,
			want:  "nil",
		},
		"string": {
			input: "test",
			want:  "string",
		},
		"std lib": {
			input: io.Discard,
			want:  "io.discard",
		},
		"external lib": {
			input: vector3.New(1, 2, 3),
			want:  "vector3.Vector[int]",
		},
		"pointer external lib": {
			input: &vector3.Vector[float64]{},
			want:  "vector3.Vector[float64]",
		},
		"nil pointer external lib": {
			input: v,
			want:  "vector3.Vector[float64]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := refutil.GetTypeName(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

type ParentGeneric[T any] struct{}

func TestGetTypeNameWithoutPackage(t *testing.T) {
	// var reader io.Reader
	var v *vector3.Vector[float64]
	tests := map[string]struct {
		input any
		want  string
	}{
		"nil": {
			input: nil,
			want:  "nil",
		},
		"string": {
			input: "test",
			want:  "string",
		},
		"std lib": {
			input: io.Discard,
			want:  "discard",
		},
		"external lib": {
			input: vector3.New(1, 2, 3),
			want:  "Vector[int]",
		},
		"pointer external lib": {
			input: &vector3.Vector[float64]{},
			want:  "Vector[float64]",
		},
		"nil pointer external lib": {
			input: v,
			want:  "Vector[float64]",
		},
		"nested generics": {
			input: ParentGeneric[vector3.Vector[float64]]{},
			want:  "ParentGeneric[github.com/EliCDavis/vector/vector3.Vector[float64]]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := refutil.GetTypeNameWithoutPackage(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSetStructField(t *testing.T) {
	v := vector3.New(1., 2., 3.)
	ts := &TestStruct{
		A: 6,
		B: true,
		C: vector3.New(1., 2., 3.),
		D: &v,
	}

	refutil.SetStructField(ts, "A", 4)
	refutil.SetStructField(ts, "B", false)
	refutil.SetStructField(ts, "C", vector3.New(4., 5., 6.))
	refutil.SetStructField(ts, "D", nil)

	assert.Equal(t, 4, ts.A)
	assert.Equal(t, false, ts.B)
	assert.Equal(t, vector3.New(4., 5., 6.), ts.C)
	assert.Nil(t, ts.D)

	assert.PanicsWithError(t, "field 'FAIL' was not found on struct", func() {
		refutil.SetStructField(ts, "FAIL", 5)
	})

	assert.PanicsWithError(t, "value of type: 'int' has no field 'D' to set", func() {
		refutil.SetStructField(12, "D", 5)
	})
}

func TestAddToStructFieldArray(t *testing.T) {
	type Struct struct {
		Blah string
		A    []int
	}

	s := &Struct{}
	refutil.AddToStructFieldArray(s, "A", 1)
	assert.Len(t, s.A, 1)
	assert.Equal(t, 1, s.A[0])
}

func TestRemoveFromStructFieldArray(t *testing.T) {
	type Struct struct {
		Blah string
		A    []int
	}

	s := &Struct{
		A: []int{1, 2, 3},
	}
	refutil.RemoveFromStructFieldArray(s, "A", 1)
	assert.Len(t, s.A, 2)
	assert.Equal(t, 1, s.A[0])
	assert.Equal(t, 3, s.A[1])
}

func TestGetFieldValues(t *testing.T) {
	ts := TestStruct{}

	v := refutil.StructFieldTypes(ts)
	assert.Len(t, v, 4)
	assert.Equal(t, "int", v["A"])
	assert.Equal(t, "bool", v["B"])
	assert.Equal(t, "vector3.Vector[float64]", v["C"])
	assert.Equal(t, "*vector3.Vector[float64]", v["D"])

	// A int
	// B bool
	// C vector3.Float64
	// D *vector3.Float64
}
