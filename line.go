package main

import (
	"math"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func NewLine(shape *DynamicShape) *Line {
	return &Line{
		DynamicShape: shape,

		rate:   time.Second / 5.0,
		height: 1.0,
		length: 0.5,
	}
}

type Line struct {
	*DynamicShape

	lastMove time.Duration
	rate     time.Duration

	lastTurn mgl.Vec3
	position mgl.Vec3
	height   float32
	length   float32

	angle   float32
	turning bool
	offset  int
}

func (line *Line) Tick(since time.Duration, rotate float32) {
	if rotate != 0 {
		line.angle += rotate
		line.turning = true
	}

	if line.lastMove+line.rate > since {
		return
	}

	line.lastMove = since

	line.Add(line.angle)
	line.Buffer(line.offset)
}

func (line *Line) Add(angle float32) {
	turning := line.turning || line.angle != angle
	if turning {
		line.angle = angle
		line.lastTurn = line.position
	}

	// Rotate
	sin, cos := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))
	unit := mgl.Vec3{cos - sin, 0, sin + cos}

	// Normalize and reset height
	l := line.length / unit.Len()
	unit = mgl.Vec3{unit[0] * l, line.height, unit[2] * l}

	p1 := line.lastTurn
	p2 := line.position.Add(unit)
	p3 := mgl.Vec3{p2[0], 0, p2[2]} // Discard height
	quad := Quad(p1, p2)

	pn := p1.Sub(p2).Cross(p3.Sub(p2)).Normalize()
	normal := pn[:]

	line.position = p3

	if !turning && len(line.vertices) >= len(quad) {
		// Replace
		line.vertices = append(line.vertices[:len(line.vertices)-len(quad)], quad...)
	} else {
		line.offset = len(line.vertices) / vertexDim
		line.vertices = append(line.vertices, quad...)
	}
	// TODO: Optimize by using indices
	line.normals = append(line.normals, normal...)
	line.normals = append(line.normals, normal...)
	line.normals = append(line.normals, normal...)
	line.normals = append(line.normals, normal...)
	line.normals = append(line.normals, normal...)
	line.normals = append(line.normals, normal...)

	line.turning = false
}
