package main

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
)

func NewLine(shader Shader, bufSize int) *Line {
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferInit(gl.ARRAY_BUFFER, bufSize, gl.DYNAMIC_DRAW)

	position := mgl.Vec3{}

	return &Line{
		shader:    shader,
		VBO:       vbo,
		rate:      time.Second / 6,
		height:    1.0,
		step:      0.5,
		direction: mgl.Vec3{1, 0, 1}, // angle=0
		position:  position,
		segments:  []mgl.Vec3{position},
	}
}

type Line struct {
	shader Shader
	VBO    gl.Buffer

	interval time.Duration
	rate     time.Duration

	position  mgl.Vec3
	direction mgl.Vec3
	segments  []mgl.Vec3
	height    float32
	step      float32

	angle   float64
	turning bool
	offset  int
}

func (line *Line) Tick(interval time.Duration, rotate float64) {
	if rotate != 0 {
		line.angle += rotate
		line.turning = true
	}

	line.interval += interval
	if line.interval < line.rate {
		return
	}
	line.interval -= line.rate

	line.Add(line.angle)
	line.Buffer(line.offset)
}

func (line *Line) Add(angle float64) {
	turning := line.turning || line.angle != angle
	if turning {
		line.angle = angle
		sin, cos := math.Sin(angle), math.Cos(line.angle)
		line.direction = mgl.Vec3{float32(cos - sin), 0, float32(sin + cos)}
	}

	// Normalize and reset height
	unit := line.direction
	l := line.step / unit.Len()
	unit = mgl.Vec3{unit[0] * l, 0.0, unit[2] * l}
	line.position = line.position.Add(unit)

	if !turning && len(line.segments) > 1 {
		// Replace
		line.segments[len(line.segments)-1] = line.position
	} else {
		line.offset = len(line.segments)
		line.segments = append(line.segments, line.position)
	}

	line.turning = false
}

// Shape interface:

func (shape *Line) Len() int {
	return len(shape.segments) * lineEmittedVertices
}

func (shape *Line) Stride() int {
	return vecSize * vertexDim
}

func (shape *Line) Close() error {
	gl.DeleteBuffer(shape.VBO)
	return nil
}

const lineEmittedVertices = 2

func (shape *Line) BytesOffset(n int) []byte {
	quad := [6]float32{}
	buf := bytes.Buffer{}

	var s mgl.Vec3
	var bot, top float32 = 0.0, shape.height

	for i := n; i < len(shape.segments); i++ {
		s = shape.segments[i]

		quad = [6]float32{
			s[0], bot, s[2], // Bottom Right
			s[0], top, s[2], // Top Right
		}
		binary.Write(&buf, binary.LittleEndian, quad)
	}
	return buf.Bytes()
}

func (shape *Line) Buffer(offset int) {
	data := shape.BytesOffset(offset)
	if len(data) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, lineEmittedVertices*offset*shape.Stride(), data)
}

func (shape *Line) Draw(camera Camera) {
	shader := shape.shader
	gl.Uniform1f(shader.Uniform("lights[0].intensity"), 2.0)
	gl.Uniform3fv(shader.Uniform("lights[0].position"), shape.position[:])

	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
	stride := shape.Stride()

	gl.EnableVertexAttribArray(shader.Attrib("vertCoord"))
	gl.VertexAttribPointer(shader.Attrib("vertCoord"), vertexDim, gl.FLOAT, false, stride, 0)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, shape.Len())
}

// Node interface:

func (node *Line) UseShader(parent Shader) (Shader, bool) {
	if parent == node.shader {
		return parent, false
	}
	node.shader.Use()
	return node.shader, true
}

func (node *Line) Transform(parent *mgl.Mat4) mgl.Mat4 {
	return MultiMul(parent)
}
