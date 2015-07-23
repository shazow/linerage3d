package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func TestEncodeObjects(t *testing.T) {
	vertices := []float32{42}
	bytes := EncodeObjects(0, 1, NewObjectData(vertices, 1))
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
	bytes = EncodeObjects(0, 1, NewObjectData(vertices, 2))
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
	bytes = EncodeObjects(0, 2, NewObjectData(vertices, 1))
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
	objects := []ObjectData{NewObjectData(vertices, dim)}
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
	expectIdx := []int{1, 2}
	if !reflect.DeepEqual(idx, expectIdx) {
		t.Error("got %q; want %q", idx, expectIdx)
	}
	expectVerts := []float32{1, 1, 1, 2, 2, 2}
	if !reflect.DeepEqual(verts, expectVerts) {
		t.Error("got %q; want %q", verts, expectVerts)
	}

	idx = []int{}
	verts = AppendIndexed([]float32{}, &idx, skyboxVertices...)
	fmt.Println(idx, verts)
}
