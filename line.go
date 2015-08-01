package main

import (
	"math"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
)

func NewLine(shader Shader, bufSize int) *Line {
	shape := NewDynamicShape(bufSize)

	return &Line{
		Node: Node{
			Shape:  shape,
			shader: shader,
		},
		shape:     shape,
		rate:      time.Second / 6,
		height:    1.0,
		step:      0.5,
		direction: mgl.Vec3{1, 0, 1}, // angle=0
	}
}

type Line struct {
	Node
	shape *DynamicShape

	interval time.Duration
	rate     time.Duration

	lastTurn  mgl.Vec3
	position  mgl.Vec3
	direction mgl.Vec3
	height    float32
	step      float32

	angle   float64
	turning bool
	offset  int
}

func (line *Line) Draw(camera Camera) {
	shader := line.Node.shader
	gl.Uniform1f(shader.Uniform("lights[0].intensity"), 2.0)
	gl.Uniform3fv(shader.Uniform("lights[0].position"), line.position[:])

	line.Node.Draw(camera)
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
	line.shape.Buffer(line.offset)
}

func (line *Line) Add(angle float64) {
	turning := line.turning || line.angle != angle
	if turning {
		line.angle = angle
		line.lastTurn = line.position

		sin, cos := math.Sin(angle), math.Cos(line.angle)
		line.direction = mgl.Vec3{float32(cos - sin), 0, float32(sin + cos)}
	}

	// Normalize and reset height
	unit := line.direction
	l := line.step / unit.Len()
	unit = mgl.Vec3{unit[0] * l, line.height, unit[2] * l}

	p1 := line.lastTurn
	p2 := line.position.Add(unit)
	p3 := mgl.Vec3{p2[0], 0, p2[2]} // Discard height
	quad := Quad(p1, p2)

	pn := p1.Sub(p2).Cross(p3.Sub(p2)).Normalize()
	normal := pn[:]

	line.position = p3
	shape := line.shape

	if !turning && len(shape.vertices) >= len(quad) {
		// Replace
		shape.vertices = append(shape.vertices[:len(shape.vertices)-len(quad)], quad...)
	} else {
		line.offset = len(shape.vertices) / vertexDim
		shape.vertices = append(shape.vertices, quad...)
	}
	// TODO: Optimize by using indices
	shape.normals = append(shape.normals, normal...)
	shape.normals = append(shape.normals, normal...)
	shape.normals = append(shape.normals, normal...)
	shape.normals = append(shape.normals, normal...)
	shape.normals = append(shape.normals, normal...)
	shape.normals = append(shape.normals, normal...)

	line.turning = false
}
