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

type Shape struct {
	VBO   gl.Buffer
	glTex gl.Texture

	vertices []float32 // Vec3
	textures []float32 // Vec2 (UV)
	normals  []float32 // Vec3
}

func (s *Shape) Stride() int {
	r := vertexDim
	if len(s.textures) > 0 {
		r += textureDim
	}
	if len(s.normals) > 0 {
		r += normalDim
	}
	return r * vecSize
}

func (s *Shape) Len() int {
	return len(s.vertices) / vertexDim
}

func (s *Shape) BytesOffset(n int) []byte {
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

func (s *Shape) Bytes() []byte {
	return s.BytesOffset(0)
}

func (s *Shape) BufferSub(offset int) {
	gl.BindBuffer(gl.ARRAY_BUFFER, s.VBO)
	data := s.BytesOffset(offset)
	if len(data) == 0 {
		return
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, offset*s.Stride(), data)
}

func (s *Shape) Buffer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, s.VBO)
	data := s.Bytes()
	if len(data) == 0 {
		return
	}
	gl.BufferData(gl.ARRAY_BUFFER, data, gl.STATIC_DRAW)
}

// TODO: Good render loop: http://www.java-gaming.org/index.php?topic=18710.0
