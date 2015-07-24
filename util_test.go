package main

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
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

	s = NewDimSlice(2, []int{1, 2, 3, 4, 5, 6})

	if a, b := s.Dim(), 2; a != b {
		t.Error("got %q; want %q", a, b)
	}
	if a, b := s.Slice(1, 4), []int{2, 3, 4}; !reflect.DeepEqual(a, b) {
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
