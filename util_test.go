package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func DecodeObjects(b []byte) ([]float32, error) {
	size := 4
	r := []float32{}
	buf := bytes.NewReader(b)
	for i := 0; i < len(b); i += size {
		var v float32
		err := binary.Read(buf, binary.LittleEndian, &v)
		if err != nil {
			return r, err
		}
		r = append(r, v)
	}

	return r, nil
}

func TestDimSlice(t *testing.T) {
	s := NewDimSlice(3, []float32{1, 2, 3, 4, 5, 6})

	if a, b := s.Dim(), 3; a != b {
		t.Error("got %q; want %q", a, b)
	}
	if a, b := s.Slice(1, 4), []float32{2, 3, 4}; !reflect.DeepEqual(a, b) {
		t.Error("got %q; want %q", a, b)
	}

	s = NewDimSlice(2, []uint8{1, 2, 3, 4, 5, 6})

	if a, b := s.Dim(), 2; a != b {
		t.Error("got %q; want %q", a, b)
	}
	if a, b := s.Slice(1, 4), []uint8{2, 3, 4}; !reflect.DeepEqual(a, b) {
		t.Error("got %q; want %q", a, b)
	}
}

func TestEncodeObjects(t *testing.T) {
	vertices := []float32{42}
	bytes := EncodeObjects(0, 1, NewDimSlice(1, vertices))
	if len(bytes) != 4 {
		t.Error("encoded float32 slice is the wrong size:", len(bytes), "!=", 4)
	}
	decoded, err := DecodeObjects(bytes)
	if err != nil {
		t.Error("Failed to decode:", err)
	}
	if !reflect.DeepEqual([]float32{42}, decoded) {
		t.Error("Failed to encode:", decoded)
	}

	vertices = []float32{42, 123}
	bytes = EncodeObjects(0, 1, NewDimSlice(2, vertices))
	if len(bytes) != 8 {
		t.Error("encoded float32 slice is the wrong size:", len(bytes), "!=", 8)
	}
	decoded, err = DecodeObjects(bytes)
	if err != nil {
		t.Error("Failed to decode:", err)
	}
	if !reflect.DeepEqual(vertices, decoded) {
		t.Error("Failed to encode:", decoded)
	}

	vertices = []float32{42, 123}
	bytes = EncodeObjects(0, 2, NewDimSlice(1, vertices))
	if len(bytes) != 8 {
		t.Error("encoded float32 slice is the wrong size:", len(bytes), "!=", 8)
	}
	decoded, err = DecodeObjects(bytes)
	if err != nil {
		t.Error("Failed to decode:", err)
	}
	if !reflect.DeepEqual(vertices, decoded) {
		t.Error("Failed to encode:", decoded)
	}

	dim := 3
	vertices = []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
	objects := []DimSlicer{NewDimSlice(dim, vertices)}
	bytes = EncodeObjects(0, len(vertices)/dim, objects...)
	if len(bytes) != 4*len(vertices) {
		t.Error("encoded float32 slice is the wrong size:", len(bytes), "!=", 4*len(vertices))
	}

	decoded, err = DecodeObjects(bytes)
	if err != nil {
		t.Error("Failed to decode:", err)
	}
	if !reflect.DeepEqual(vertices, decoded) {
		t.Error("Failed to encode:", decoded)
	}
}

func TestQuad(t *testing.T) {
	t.SkipNow()
	q := Quad(mgl.Vec3{0, 0, 0}, mgl.Vec3{1, 1, 0})
	fmt.Println(q)

	q = Quad(mgl.Vec3{-1, -1, 0}, mgl.Vec3{1, 1, 0})
	fmt.Println(q)

	q = Quad(mgl.Vec3{0, 0, 0}, mgl.Vec3{1, 0, 1})
	fmt.Println(q)
}

func TestAppendIndexed(t *testing.T) {
	idx := []int{}
	verts := AppendIndexed([]float32{}, &idx, []float32{1, 1, 1, 2, 2, 2, 1, 1, 1, 2, 2, 2, 1, 1, 1}...)
	expectIdx := []int{0, 1, 0, 1, 0}
	if !reflect.DeepEqual(idx, expectIdx) {
		t.Errorf("got %v; want %v", idx, expectIdx)
	}
	expectVerts := []float32{1, 1, 1, 2, 2, 2}
	if !reflect.DeepEqual(verts, expectVerts) {
		t.Errorf("got %v; want %v", verts, expectVerts)
	}

	idx = []int{}
	verts = AppendIndexed([]float32{}, &idx, unindexedCube...)
	if !reflect.DeepEqual(verts, indexedCube) {
		t.Errorf("got %v; want %v", verts, indexedCube)
	}

	if !reflect.DeepEqual(idx, cubeIndex) {
		t.Errorf("got %v; want %v", idx, cubeIndex)
	}
}

var unindexedCube = []float32{
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
}

var indexedCube = []float32{
	-1, 1, -1,
	-1, -1, -1,
	1, -1, -1,
	1, 1, -1,
	-1, -1, 1,
	-1, 1, 1,
	1, -1, 1,
	1, 1, 1,
}

var cubeIndex = []int{
	0, 1, 2, 2, 3, 0,
	4, 1, 0, 0, 5, 4,
	2, 6, 7, 7, 3, 2,
	4, 5, 7, 7, 6, 4,
	0, 3, 7, 7, 5, 0,
	1, 4, 2, 2, 4, 6,
}

func TestBoxCollision(t *testing.T) {
	tests := []struct {
		result bool

		a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y float32
	}{
		{true, 0, 0, 1, 1, 0, 0, 1, 1},
		{true, 0, 0, 1, 1, 1, 1, 0, 0},
		{true, 0, 0, 1, 1, 1, 0, 0, 1},
		{false, 0, 0, 1, 1, 2, 2, 3, 3},
	}

	for i, test := range tests {
		r := IsBoxCollision(test.a1_x, test.a1_y, test.a2_x, test.a2_y, test.b1_x, test.b1_y, test.b2_x, test.b2_y)
		if r != test.result {
			t.Errorf("IsBoxCollision test #%d failed: %v", i, test)
		}
	}
}

func TestIsCollision(t *testing.T) {
	tests := []struct {
		result bool

		a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y float32
	}{
		// Test collision is [a, b)
		{true, 0, 0, 1, 1, 1, 1, 0, 0},      // a -> -a
		{true, 0, 0, 1, 1, 0, 0, 1, 1},      // a -> a
		{true, 0, 0, 1, 1, 1, 1, 0.5, 0.5},  // a -> b intersect
		{false, 0, 0, 1, 1, 1, 1, 2, 2},     // a -> b
		{true, 2, 2, 1, 1, 1, 1, 1.5, 1.5},  // b -> a intersect
		{false, 2, 2, 1, 1, 1, 1, 0.5, 0.5}, // b -> a
		{false, 1, 0, 2, 0, 2, 0, 3, 0},     // a -> b
		{false, 0, 0, 1, 1, 1, 1, 1.5, 1.5}, // a -> b, non-collinear
		{true, 0, 0, 4, 4, 1, 1, 2, 2},      // a -> d, b -> c, collinear contained
		{true, 0, 0, 4, 4, 2, 2, 1, 1},      // a -> d, c -> b, collinear contained
		{false, 0, 1, 0, 4, 1, 3, 1, 2},     // a -> d, b -> c, parallel offset
		{false, 0, 1, 0, 4, 1, 2, 1, 3},     // a -> d, c -> b, parallel offset
		{true, 1, 0, 2, 0, 3, 0, 2, 0},      // a -> b <- c

		{true, 0, 1, 0, 4, 0, 1, 1, 3}, // a -> b, a->c, non-collinear
		{true, 0, 1, 0, 4, 0, 1, 0, 3},
		{false, 0, 1, 0, 4, 0, 4, 1, 3},
		{true, 0, 1, 0, 4, 0, 4, 0, 3},

		{true, 3, 1, 3, 2, 2.5, 1.5, 3.5, 1.5},
		{true, 3, 1, 3, 2, 2.5, 1.5, 3.5, 1},
		{true, 3, 1, 3, 2, 2.5, 1, 4, 1},
		{true, 1, 1, 1, 3, 0, 1, 4, 1},
		{true, 1, 0, 1, 1, 0, 1, 4, 1},
		{false, 0, 0, 1, 1, 2, 2, 3, 3},
		{false, 2, 0, 3, 0, 3, 0, 3, 1},
		{false, 2, 1, 3, 1, 1, 0, 2, 0}, // collinear disjoint vertically
		//{false, 1.37, 1.39, 1.34, 1.35, 1.31, 1.31, 1.34, 1.35},
	}

	for i, test := range tests {
		r := IsCollision2D(test.a1_x, test.a1_y, test.a2_x, test.a2_y, test.b1_x, test.b1_y, test.b2_x, test.b2_y)
		if r != test.result {
			t.Errorf("IsCollision2D test #%d failed: %v", i, test)
		}
	}
}
