package main

import "golang.org/x/mobile/gl"

const vertexDim = 3
const textureDim = 2
const normalDim = 3
const vecSize = 4

type Shape interface {
	Close() error
	Stride() int
	Len() int
	Draw(Shader, Camera)
}

func NewStaticShape() *StaticShape {
	shape := &StaticShape{}
	shape.VBO = gl.CreateBuffer()
	shape.IBO = gl.CreateBuffer()
	return shape
}

type StaticShape struct {
	VBO     gl.Buffer
	IBO     gl.Buffer
	Texture gl.Texture

	vertices []float32 // Vec3
	textures []float32 // Vec2 (UV)
	normals  []float32 // Vec3
	indices  []uint8
}

func (s *StaticShape) Len() int {
	return len(s.vertices) / vertexDim
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

func (shape *StaticShape) Bytes() []byte {
	return shape.BytesOffset(0)
}

func (shape *StaticShape) Close() error {
	gl.DeleteBuffer(shape.VBO)
	gl.DeleteBuffer(shape.IBO)
	return nil
}

func (shape *StaticShape) BytesOffset(n int) []byte {
	objects := []DimSlicer{NewDimSlice(vertexDim, shape.vertices)}
	if len(shape.textures) > 0 {
		objects = append(objects, NewDimSlice(textureDim, shape.textures))
	}
	if len(shape.normals) > 0 {
		objects = append(objects, NewDimSlice(normalDim, shape.normals))
	}

	length := len(shape.vertices) / vertexDim
	return EncodeObjects(n, length, objects...)
}

func (shape *StaticShape) Draw(shader Shader, camera Camera) {
	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
	stride := shape.Stride()

	gl.EnableVertexAttribArray(shader.Attrib("vertCoord"))
	gl.VertexAttribPointer(shader.Attrib("vertCoord"), vertexDim, gl.FLOAT, false, stride, 0)

	if len(shape.normals) > 0 {
		gl.EnableVertexAttribArray(shader.Attrib("vertNormal"))
		gl.VertexAttribPointer(shader.Attrib("vertNormal"), normalDim, gl.FLOAT, false, stride, vertexDim*vecSize)
	}
	// TODO: texture

	if len(shape.indices) > 0 {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, shape.IBO)
		gl.DrawElements(gl.TRIANGLES, len(shape.indices), gl.UNSIGNED_BYTE, 0)
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, shape.Len())
	}

	gl.DisableVertexAttribArray(shader.Attrib("vertCoord"))
	if len(shape.normals) > 0 {
		gl.DisableVertexAttribArray(shader.Attrib("vertNormal"))
	}
	// TODO: texture
}

func (shape *StaticShape) Buffer() {
	data := shape.Bytes()
	if len(data) > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
		gl.BufferData(gl.ARRAY_BUFFER, data, gl.STATIC_DRAW)
	}

	if len(shape.indices) > 0 {
		data = EncodeObjects(0, len(shape.indices), NewDimSlice(1, shape.indices))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, shape.IBO)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, data, gl.STATIC_DRAW)
	}
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
