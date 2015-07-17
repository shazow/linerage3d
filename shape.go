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
	Bind()
	Stride() int
	Len() int
	EnableAttrib(gl.Attrib, gl.Attrib, gl.Attrib)
	DisableAttrib(gl.Attrib, gl.Attrib, gl.Attrib)
}

type DynamicShape struct {
	VBO   gl.Buffer
	glTex gl.Texture

	vertices []float32 // Vec3
	textures []float32 // Vec2 (UV)
	normals  []float32 // Vec3
}

func (shape *DynamicShape) EnableAttrib(vertex gl.Attrib, normal gl.Attrib, texture gl.Attrib) {
	shape.Bind()
	stride := shape.Stride()

	gl.EnableVertexAttribArray(vertex)
	gl.VertexAttribPointer(vertex, vertexDim, gl.FLOAT, false, stride, 0)

	if len(shape.normals) > 0 {
		gl.EnableVertexAttribArray(normal)
		gl.VertexAttribPointer(normal, normalDim, gl.FLOAT, false, stride, vertexDim*vecSize)
	}
	// TODO: texture
}

func (shape *DynamicShape) DisableAttrib(vertex gl.Attrib, normal gl.Attrib, texture gl.Attrib) {
	gl.DisableVertexAttribArray(vertex)
	if len(shape.normals) > 0 {
		gl.DisableVertexAttribArray(normal)
	}
	// TODO: texture
}

func (s *DynamicShape) Stride() int {
	r := vertexDim
	if len(s.textures) > 0 {
		r += textureDim
	}
	if len(s.normals) > 0 {
		r += normalDim
	}
	return r * vecSize
}

func (s *DynamicShape) Len() int {
	return len(s.vertices) / vertexDim
}

func (s *DynamicShape) BytesOffset(n int) []byte {
	buf := bytes.Buffer{}

	wrote := [][]float32{}

	for i := n; i < s.Len(); i++ {
		v := s.vertices[i*vertexDim : i*vertexDim+3]
		if len(s.textures) > 0 {
			v = append(v, s.textures[i*textureDim:(i+1)*textureDim]...)
		}
		if len(s.normals) > 0 {
			v = append(v, s.normals[i*normalDim:(i+1)*normalDim]...)
		}

		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			panic(fmt.Sprintln("binary.Write failed:", err))
		}

		wrote = append(wrote, v)
	}

	fmt.Printf("Wrote %d vertices: %d to %d \t", s.Len()-n, n, s.Len())
	fmt.Println(wrote)

	return buf.Bytes()
}

func (s *DynamicShape) Bytes() []byte {
	return s.BytesOffset(0)
}

// Bind activates the shape's buffers
func (s *DynamicShape) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, s.VBO)
}

func (s *DynamicShape) Init(bufSize int) {
	s.VBO = gl.CreateBuffer()
	s.Bind()
	gl.BufferInit(gl.ARRAY_BUFFER, bufSize, gl.DYNAMIC_DRAW)
}

func (s *DynamicShape) Buffer(offset int) {
	s.Bind()
	data := s.BytesOffset(offset)
	if len(data) == 0 {
		return
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, offset*s.Stride(), data)
}

// TODO: Good render loop: http://www.java-gaming.org/index.php?topic=18710.0
