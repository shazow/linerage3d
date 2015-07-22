package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"golang.org/x/mobile/gl"
)

const vertexDim = 3
const textureDim = 2
const normalDim = 3
const vecSize = 4

type Shape interface {
	Close()
	Stride() int
	Len() int
	EnableAttrib(gl.Attrib, gl.Attrib, gl.Attrib)
	DisableAttrib(gl.Attrib, gl.Attrib, gl.Attrib)
}

func NewStaticShape() *StaticShape {
	shape := &StaticShape{}
	shape.VBO = gl.CreateBuffer()
	shape.EBI = gl.CreateBuffer()
	return shape
}

type StaticShape struct {
	VBO   gl.Buffer
	EBI   gl.Buffer
	glTex gl.Texture

	vertices []float32 // Vec3
	textures []float32 // Vec2 (UV)
	normals  []float32 // Vec3
	indexes  []float32 // Vec3
}

func (shape *StaticShape) EnableAttrib(vertex gl.Attrib, normal gl.Attrib, texture gl.Attrib) {
	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
	stride := shape.Stride()

	gl.EnableVertexAttribArray(vertex)
	gl.VertexAttribPointer(vertex, vertexDim, gl.FLOAT, false, stride, 0)

	if len(shape.normals) > 0 {
		gl.EnableVertexAttribArray(normal)
		gl.VertexAttribPointer(normal, normalDim, gl.FLOAT, false, stride, vertexDim*vecSize)
	}
	// TODO: texture
}

func (shape *StaticShape) DisableAttrib(vertex gl.Attrib, normal gl.Attrib, texture gl.Attrib) {
	gl.DisableVertexAttribArray(vertex)
	if len(shape.normals) > 0 {
		gl.DisableVertexAttribArray(normal)
	}
	// TODO: texture
}

func (shape *StaticShape) Stride() int {
	r := vertexDim
	if len(shape.textures) > 0 {
		r += textureDim
	}
	if len(shape.normals) > 0 {
		r += normalDim
	}
	return r * vecSize
}

func (s *StaticShape) Len() int {
	return len(s.vertices) / vertexDim
}

func (shape *StaticShape) BytesOffset(n int) []byte {
	buf := bytes.Buffer{}

	wrote := [][]float32{}

	for i := n; i < shape.Len(); i++ {
		v := shape.vertices[i*vertexDim : i*vertexDim+3]
		if len(shape.textures) > 0 {
			v = append(v, shape.textures[i*textureDim:(i+1)*textureDim]...)
		}
		if len(shape.normals) > 0 {
			v = append(v, shape.normals[i*normalDim:(i+1)*normalDim]...)
		}

		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			panic(fmt.Sprintln("binary.Write failed:", err))
		}

		wrote = append(wrote, v)
	}

	//fmt.Printf("Wrote %d vertices: %d to %d \t", shape.Len()-n, n, shape.Len())
	//fmt.Println(wrote)

	return buf.Bytes()
}

func (shape *StaticShape) Bytes() []byte {
	return shape.BytesOffset(0)
}

func (shape *StaticShape) Close() {
	gl.DeleteBuffer(shape.VBO)
	gl.DeleteBuffer(shape.EBI)
}

func (shape *StaticShape) Buffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
	data := shape.Bytes()
	if len(data) == 0 {
		return
	}
	gl.BufferData(gl.ARRAY_BUFFER, data, gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, shape.EBI)
}

func NewDynamicShape(bufSize int) *DynamicShape {
	shape := &DynamicShape{}
	shape.VBO = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
	gl.BufferInit(gl.ARRAY_BUFFER, bufSize, gl.DYNAMIC_DRAW)

	return shape
}

type DynamicShape struct {
	StaticShape
}

func (shape *DynamicShape) Buffer(offset int) {
	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
	data := shape.BytesOffset(offset)
	if len(data) == 0 {
		return
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, offset*shape.Stride(), data)
}

// TODO: Good render loop: http://www.java-gaming.org/index.php?topic=18710.0
